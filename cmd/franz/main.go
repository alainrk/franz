package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var broker = "localhost:29092"

type KafkaConsumer interface {
	GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error)
}

// GetAvailableTopics returns a list of all available topics in the Kafka broker.
// If excludeInternals is true, it will exclude internal topics (starting with _).
func GetAvailableTopics(client KafkaConsumer, excludeInternals bool) ([]string, error) {
	metadata, err := client.GetMetadata(nil, true, 5000)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	var topics []string
	for topic := range metadata.Topics {
		if excludeInternals && strings.HasPrefix(topic, "_") {
			continue
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "franz",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %v\n", err)
		return
	}
	defer c.Close()

	topics, err := GetAvailableTopics(c, true)
	if err != nil {
		fmt.Printf("Failed to get available topics: %v\n", err)
		return
	}

	fmt.Printf("Auto-discovered topics: %v\n", topics)

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topics: %v\n", err)
		return
	}

	fmt.Println("Consumer started. Press Ctrl+C to exit.")

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Println("Received a new message:")
				fmt.Printf("Headers: %v\n", e.Headers)
				fmt.Printf("Content: %s\n", e.Value)
				fmt.Printf("Key: %s\n", e.Key)
				fmt.Printf("Opaque: %v\n", e.Opaque)
				fmt.Printf("TopicPartition.Topic: %s\n", *e.TopicPartition.Topic)
				fmt.Printf("TopicPartition.Partition: %d\n", e.TopicPartition.Partition)
				fmt.Printf("TopicPartition.Offset: %v\n", e.TopicPartition.Offset)
				fmt.Printf("TopicPartition.Error: %v\n", e.TopicPartition.Error)
				if e.TopicPartition.Metadata != nil {
					fmt.Println("TopicPartition.Metadata:", e.TopicPartition.Metadata)
				}
				fmt.Println("-------------------------------------------")
			case kafka.PartitionEOF:
				fmt.Printf("Reached end of partition %v\n", e)
			case kafka.Error:
				fmt.Printf("Error: %v\n", e)
			}
		}
	}
}
