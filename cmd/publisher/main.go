package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"franz/models"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var broker = "localhost:29092"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <topic> <number_of_messages> [-s]")
		return
	}

	topic := os.Args[1]
	numMessages, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid number of messages:", os.Args[2])
		return
	}

	// TODO: do something properly with this flag
	srEnabled := len(os.Args) == 4 && os.Args[3] == "-s"
	var serializer *protobuf.Serializer

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	if srEnabled {
		sr, _ := schemaregistry.NewClient(schemaregistry.NewConfig("localhost:8081"))
		serializer, _ = protobuf.NewSerializer(sr, serde.ValueSerde, protobuf.NewSerializerConfig())
	}

	defer p.Close()

	fmt.Printf("Publishing %d random messages to topic '%s'...\n", numMessages, topic)

	for i := 1; i <= numMessages; i++ {
		value := getRandomMessage()
		headers := []kafka.Header{
			{Key: "message_number", Value: []byte(strconv.Itoa(i))},
			{Key: "timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		}

		err := publishBytesMessage(p, topic, []byte(strconv.Itoa(i)), value, headers)
		if err != nil {
			fmt.Printf("Failed to publish message %d: %v\n", i, err)
			break
		}

		fmt.Printf("Published message %d\n", i)
	}

	fmt.Println("All messages published!")
}

func publishBytesMessage(p *kafka.Producer, topic string, key, value []byte, headers []kafka.Header) error {
	deliveryChan := make(chan kafka.Event)

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
		Headers:        headers,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

func publishSerializedMessage(p *kafka.Producer, topic string, key, value []byte, headers []kafka.Header) error {
	deliveryChan := make(chan kafka.Event)

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
		Headers:        headers,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

func getRandomMessage() []byte {
	u := models.NewUser{
		Firstname:   "John",
		Lastname:    "Doe",
		DateOfBirth: "1990-01-01",
		HeightCm:    170 + rand.Uint32()%30,
	}
	return []byte(u.String())
}
