package main

import (
	"reflect"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// MockConsumer is a mock Kafka consumer for testing purposes.
type MockConsumer struct{}

// GetMetadata is a mock implementation of GetMetadata for MockConsumer.
func (c *MockConsumer) GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error) {
	topic1 := kafka.TopicMetadata{
		Topic:      "topic1",
		Partitions: []kafka.PartitionMetadata{1: {}},
	}
	topic2 := kafka.TopicMetadata{
		Topic:      "topic2",
		Partitions: []kafka.PartitionMetadata{2: {}},
	}
	internalTopic := kafka.TopicMetadata{
		Topic:      "__internal",
		Partitions: []kafka.PartitionMetadata{3: {}},
	}

	if allTopics {
		return &kafka.Metadata{
			Topics: map[string]kafka.TopicMetadata{
				topic1.Topic:        topic1,
				topic2.Topic:        topic2,
				internalTopic.Topic: internalTopic,
			},
		}, nil
	}

	return &kafka.Metadata{
		Topics: map[string]kafka.TopicMetadata{
			*topic: {
				Topic:      *topic,
				Partitions: []kafka.PartitionMetadata{1: {}},
			},
		},
	}, nil
}

func TestGetAvailableTopics(t *testing.T) {
	mockConsumer := &MockConsumer{}
	topics, err := GetAvailableTopics(mockConsumer, true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedTopics := []string{"topic1", "topic2"}
	if !reflect.DeepEqual(topics, expectedTopics) {
		t.Errorf("Expected topics %v, got %v", expectedTopics, topics)
	}
}

func TestGetAvailableTopicsInternalIncluded(t *testing.T) {
	mockConsumer := &MockConsumer{}
	topics, err := GetAvailableTopics(mockConsumer, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedTopics := []string{"topic1", "topic2", "__internal"}
	if !reflect.DeepEqual(topics, expectedTopics) {
		t.Errorf("Expected topics %v, got %v", expectedTopics, topics)
	}
}
