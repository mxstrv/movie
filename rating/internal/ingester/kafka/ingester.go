package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"movieapp/rating/pkg/model"
)

// Ingester defines a Kafka ingester.
type Ingester struct {
	consumer *kafka.Consumer
	topic    string
}

// NewIngester creates a new Kafka ingester.
func NewIngester(addr string, groupId string, topic string) (*Ingester, error) {
	cons, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &Ingester{cons, topic}, nil
}

// Ingest starts ingestion from Kafka and returns a channel with rating events.
func (ingester *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	if err := ingester.consumer.SubscribeTopics([]string{ingester.topic}, nil); err != nil {
		return nil, err
	}
	ch := make(chan model.RatingEvent, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				ingester.consumer.Close()
			default:
				msg, err := ingester.consumer.ReadMessage(-1)
				if err != nil {
					fmt.Printf("Consumer error: %v\n", err)
					continue
				}
				var event model.RatingEvent
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					fmt.Printf("Unmarshal error: %v\n", err)
				}
				ch <- event
			}
		}
	}()
	return ch, nil
}
