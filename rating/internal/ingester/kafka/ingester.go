package kafka

import "github.com/confluentinc/confluent-kafka-go/kafka"

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
