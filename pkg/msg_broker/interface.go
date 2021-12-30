package msgbroker

import (
	"context"

	"github.com/Shopify/sarama"
)

type MsgBroker interface {
	Publish(ctx context.Context, topic string, data string) error
	GetKafkaConn() (sarama.Client, sarama.ConsumerGroup)
}
