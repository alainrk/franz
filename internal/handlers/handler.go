package handlers

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	HandleMessage(message *kafka.Message) error
}

type DefaultHandler struct{}

func NewDefaultHandler() Handler {
	return DefaultHandler{}
}

func (h DefaultHandler) HandleMessage(message *kafka.Message) error {
	log.Infof("Received a new message from topic %s", *message.TopicPartition.Topic)

	s := fmt.Sprintf("Headers: %v\n", message.Headers)
	s += fmt.Sprintf("Content: %s\n", message.Value)
	s += fmt.Sprintf("Key: %s\n", message.Key)
	s += fmt.Sprintf("Opaque: %v\n", message.Opaque)
	s += fmt.Sprintf("TopicPartition.Topic: %s", *message.TopicPartition.Topic)
	s += fmt.Sprintf("TopicPartition.Partition: %d\n", message.TopicPartition.Partition)
	s += fmt.Sprintf("TopicPartition.Offset: %v\n", message.TopicPartition.Offset)
	s += fmt.Sprintf("TopicPartition.Error: %v\n", message.TopicPartition.Error)
	if message.TopicPartition.Metadata != nil {
		s += fmt.Sprintf("TopicPartition.Metadata: %v\n", message.TopicPartition.Metadata)
	}

	log.Infof(s)

	return nil
}
