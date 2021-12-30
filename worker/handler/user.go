package handler

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	msgbroker "github.com/mine/just-projecting/pkg/msg_broker"
	"github.com/sirupsen/logrus"
)

type UserWorkerHandler struct {
	Ready chan bool
	Ctx   context.Context
	Cfg   msgbroker.KafkaConfig
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (handler UserWorkerHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(handler.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (handler UserWorkerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (handler UserWorkerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		logrus.Infof("Message claimed: topic = %s, value = %s, timestamp = %v ", message.Topic, string(message.Value), message.Timestamp)

		switch message.Topic {
		case handler.Cfg.ExampleTopic:
			// handle message here
			fmt.Printf("handle message topic %s, here", handler.Cfg.ExampleTopic)
		default:
			// do nothing
		}

		session.MarkMessage(message, "")
	}

	return nil
}
