package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var broker = "localhost:29092"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <topic> <number_of_messages>")
		return
	}

	topic := os.Args[1]
	numMessages, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid number of messages:", os.Args[2])
		return
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	fmt.Printf("Publishing %d random messages to topic '%s'...\n", numMessages, topic)

	for i := 1; i <= numMessages; i++ {
		value := getRandomMessage()
		err := publishMessage(p, topic, []byte(strconv.Itoa(i)), value)
		if err != nil {
			fmt.Printf("Failed to publish message %d: %v\n", i, err)
			break
		}

		fmt.Printf("Published message %d\n", i)
	}

	fmt.Println("All messages published!")
}

func publishMessage(p *kafka.Producer, topic string, key, value []byte) error {
	deliveryChan := make(chan kafka.Event)

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
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
	message := fmt.Sprintf("Random Message - %d", rand.Intn(1000))
	return []byte(message)
}
