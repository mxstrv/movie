package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"movieapp/rating/pkg/model"
	"os"
	"time"
)

func main() {
	fmt.Println("Creating a Kafka producer")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	const fileName = "ratings-data.json"
	fmt.Println("Reading ratings event from file: " + fileName)

	ratingsEvents, err := readRatingEvents(fileName)
	if err != nil {
		panic(err)
	}

	var topic = "ratings"
	if err := produceRatingEvents(topic, producer, ratingsEvents); err != nil {
		panic(err)
	}

	var timeout = time.Second * 10
	fmt.Printf("Waiting %s until all events get produced\n", timeout.String())
	producer.Flush(int(timeout.Milliseconds()))

}

// readRatingEvents reads a json file and returns a slice of RatingEvent
func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

// produceRatingEvents creates a mock Produce requests to Kafka instance.
func produceRatingEvents(topic string, producer *kafka.Producer, events []model.RatingEvent) error {
	for _, e := range events {
		encoded, err := json.Marshal(e)
		if err != nil {
			return err
		}

		// kafka.PartitionAny for testing purposes
		// TODO set a strict partition
		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          encoded,
		}, nil); err != nil {
			return err
		}
	}
	return nil
}
