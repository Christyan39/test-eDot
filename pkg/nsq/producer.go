package nsq

import (
	"encoding/json"
	"fmt"

	"github.com/nsqio/go-nsq"
)

// Producer wraps the nsq.Producer
type Producer struct {
	producer *nsq.Producer
}

// NewProducer creates a new NSQ producer
func NewProducer(addr string, config *nsq.Config) (*Producer, error) {
	p, err := nsq.NewProducer(addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create NSQ producer: %w", err)
	}
	return &Producer{producer: p}, nil
}

// PublishJSON publishes a message to the given topic, marshaling the message as JSON
func (p *Producer) PublishJSON(topic string, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	return p.producer.Publish(topic, body)
}

// Stop stops the producer
func (p *Producer) Stop() {
	p.producer.Stop()
}
