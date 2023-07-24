package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var broker = "localhost:29092"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <topic1> <topic2> ...")
		return
	}

	topics := os.Args[1:]

	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatalf("Failed to create admin client: %v", err)
	}
	defer adminClient.Close()

	for _, topic := range topics {
		// Ask user for confirmation
		var confirmation string
		fmt.Printf("Are you sure you want to permanently delete topic '%s'?\n", topic)
		fmt.Printf("Type 'yes' to confirm: ")
		fmt.Scanln(&confirmation)
		if confirmation != "yes" {
			fmt.Println("Aborting...")
			return
		}

		fmt.Printf("Deleting topic '%s'...\n", topic)

		err := deleteTopic(adminClient, topic)
		if err != nil {
			fmt.Printf("Failed to delete topic '%s': %v\n", topic, err)
		} else {
			fmt.Printf("Topic '%s' permanently deleted successfully!\n", topic)
		}
	}
}

func deleteTopic(adminClient *kafka.AdminClient, topic string) error {
	// Get topic metadata to ensure the topic exists
	meta, err := adminClient.GetMetadata(&topic, false, 5000)
	if err != nil {
		return fmt.Errorf("failed to get topic metadata: %w", err)
	}

	if len(meta.Topics) == 0 {
		return fmt.Errorf("topic '%s' does not exist", topic)
	}

	// Delete the topic
	results, err := adminClient.DeleteTopics(context.Background(), []string{topic}, kafka.SetAdminOperationTimeout(5000))
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	// Check the results
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			return fmt.Errorf("failed to delete topic: %w", result.Error)
		}
	}

	return nil
}
