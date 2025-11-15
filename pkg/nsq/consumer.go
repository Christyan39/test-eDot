package nsq

import (
	"github.com/nsqio/go-nsq"
)

// MessageHandler is a generic handler for processing NSQ messages
type MessageHandler struct {
	ProcessFunc func(msg *nsq.Message) error
}

func (h *MessageHandler) HandleMessage(msg *nsq.Message) error {
	if h.ProcessFunc != nil {
		return h.ProcessFunc(msg)
	}
	return nil
}

// NewConsumer creates and configures a generic NSQ consumer
func NewConsumer(topic, channel, nsqdAddr string, config *nsq.Config, handler nsq.Handler) (*nsq.Consumer, error) {
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}
	consumer.AddHandler(handler)
	if err := consumer.ConnectToNSQD(nsqdAddr); err != nil {
		return nil, err
	}
	return consumer, nil
}
