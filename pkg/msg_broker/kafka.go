package msgbroker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"

	"github.com/mine/just-projecting/pkg/config"
	"github.com/mine/just-projecting/pkg/logger"
)

type kafkaMsg struct {
	logger        *logrus.Entry
	producer      sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
	cli           sarama.Client
}

func NewKafkaMsg(ctx context.Context) (MsgBroker, error) {
	log := logger.GetLoggerContext(ctx, "msgBroker", "NewKafka")

	jsonByte, err := json.Marshal(config.Get("kafka"))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var cfg KafkaConfig
	err = json.Unmarshal(jsonByte, &cfg)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	kafkaConfig := getKafkaConfig(cfg)
	hostPort := getURLKafka(cfg.UrlKafkaList)

	producer, err := sarama.NewSyncProducer(hostPort, kafkaConfig)
	if err != nil {
		log.Fatalf("Failed Create new producer, Error : %v", err)
		return nil, err
	}

	consGroup, err := sarama.NewConsumerGroup(hostPort, cfg.ConsumerGroup, kafkaConfig)
	if err != nil {
		log.Fatalf("Failed create consumer group. Error : %v", err)
		return nil, err
	}

	client, err := sarama.NewClient(hostPort, kafkaConfig)
	if err != nil {
		log.Fatalf("Failed create new kafka client. Error : %v", err)
		return nil, err
	}

	// return producer, consGroup, client
	return &kafkaMsg{
		logger:        logger.GetLogger("msgBroker", "newKafka"),
		producer:      producer,
		consumerGroup: consGroup,
		cli:           client,
	}, nil
}

func getURLKafka(urlList string) []string {
	return strings.Split(urlList, ";")
}

func getKafkaConfig(cfg KafkaConfig) *sarama.Config {
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = cfg.ProducerReturnSuccess
	kafkaConfig.Net.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Second
	kafkaConfig.Producer.Retry.Max = cfg.MaxRetry
	kafkaConfig.Version = version
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	if cfg.Username != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = cfg.Username
		kafkaConfig.Net.SASL.Password = cfg.Password
	}
	return kafkaConfig
}

// Publish publish message to kafka, return error if failed to publish
func (k *kafkaMsg) Publish(ctx context.Context, topic string, data string) error {
	log := logger.GetLoggerContext(ctx, "msgBroker", "publishKafka")
	p, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(data),
		Offset:    -1,
		Timestamp: time.Now(),
	})
	log.Infof(fmt.Sprintf("Publish Kafka Message for Topic : %s, Partition : %d, error : %v", topic, p, err))
	return err
}

func (k *kafkaMsg) GetKafkaConn() (sarama.Client, sarama.ConsumerGroup) {
	return k.cli, k.consumerGroup
}
