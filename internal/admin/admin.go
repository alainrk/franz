package admin

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaAdminer interface {
	GetTopics(includeInternals bool) ([]string, error)
}

type KafkaConsumer interface {
	GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error)
}

type DefaultKafkaAdminer struct {
	Client KafkaConsumer
}

func NewAdminClient(c KafkaConsumer) (KafkaAdminer, error) {
	return &DefaultKafkaAdminer{
		Client: c,
	}, nil
}

// GetTopics returns a list of all available topics in the Kafka broker.
// If includeInternals is true, it will include internal topics (starting with _).
func (c *DefaultKafkaAdminer) GetTopics(includeInternals bool) ([]string, error) {
	metadata, err := c.Client.GetMetadata(nil, true, 5000)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	var topics []string
	for topic := range metadata.Topics {
		if !includeInternals && strings.HasPrefix(topic, "_") {
			continue
		}
		topics = append(topics, topic)
	}

	return topics, nil
}
