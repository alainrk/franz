package main

import (
	"fmt"
	"franz/internal/admin"
	"franz/internal/config"
	"franz/internal/handlers"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var broker = "localhost:29092"

func main() {
	// Set up ctrl-c handler
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Load configuration
	config, err := config.NewConfiguration()
	fmt.Println("URLs:", config.KafkaBootstrapURLs)

	// Setup global logger
	log.SetFormatter(&log.JSONFormatter{})
	l := os.Getenv("LOG_TEXT")
	if l != "" {
		log.SetFormatter(&log.TextFormatter{})
	}
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// Create a new consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               config.KafkaBootstrapURLs,
		"group.id":                        config.KafkaConsumerGroup,
		"go.application.rebalance.enable": true,
		"session.timeout.ms":              6000,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	})
	if err != nil {
		log.Errorf("Failed to create consumer: %v", err)
		return
	}
	defer c.Close()

	// Create a new kafka admin client
	adminClient, err := admin.NewAdminClient(c)
	if err != nil {
		log.Errorf("Failed to create admin client: %v", err)
		return
	}

	// Retrieve all topics, excluding internal ones
	topics, err := adminClient.GetTopics(false)
	if err != nil {
		log.Errorf("Failed to get available topics: %v", err)
		return
	}

	log.Infof("Auto-discovered topics: %v", topics)

	// Subscribe all non-internal topics
	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		log.Errorf("Failed to subscribe to topics: %v", err)
		return
	}

	// Create a new handler
	handler := handlers.NewDefaultHandler()

	log.Info("Consumer started. Press Ctrl+C to exit.")

	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Infof("Caught signal %v: terminating", sig)
			c.Close()
			run = false
		default:
			ev := c.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				err := handler.HandleMessage(e)
				if err != nil {
					log.Errorf("Failed to handle message: %v", err)
				}
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Kafka Error on consume: %v\n", e)
				run = false
			default:
				if e != nil {
					fmt.Printf("Ignored %v\n", e)
				}
				continue
			}
		}
	}
}
