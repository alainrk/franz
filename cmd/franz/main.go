package main

import (
	"fmt"
	"franz/internal/admin"
	"franz/internal/config"
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

	log.SetFormatter(&log.JSONFormatter{})
	l := os.Getenv("LOG_TEXT")
	if l != "" {
		log.SetFormatter(&log.TextFormatter{})
	}
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

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

	adminClient, err := admin.NewAdminClient(c)
	if err != nil {
		log.Errorf("Failed to create admin client: %v", err)
		return
	}

	topics, err := adminClient.GetTopics(false)
	if err != nil {
		log.Errorf("Failed to get available topics: %v", err)
		return
	}

	log.Infof("Auto-discovered topics: %v", topics)

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		log.Errorf("Failed to subscribe to topics: %v", err)
		return
	}

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
				log.Infof("Received a new message from topic %s", *e.TopicPartition.Topic)

				s := fmt.Sprintf("Headers: %v", e.Headers)
				s += fmt.Sprintf("Content: %s", e.Value)
				s += fmt.Sprintf("Key: %s", e.Key)
				s += fmt.Sprintf("Opaque: %v", e.Opaque)
				s += fmt.Sprintf("TopicPartition.Topic: %s", *e.TopicPartition.Topic)
				s += fmt.Sprintf("TopicPartition.Partition: %d", e.TopicPartition.Partition)
				s += fmt.Sprintf("TopicPartition.Offset: %v", e.TopicPartition.Offset)
				s += fmt.Sprintf("TopicPartition.Error: %v", e.TopicPartition.Error)
				if e.TopicPartition.Metadata != nil {
					s += fmt.Sprintf("TopicPartition.Metadata: %v", e.TopicPartition.Metadata)
				}
				log.Infof(s)

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
