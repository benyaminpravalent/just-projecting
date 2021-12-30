package worker

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/mine/just-projecting/api/helper"
	"github.com/mine/just-projecting/pkg/config"
	"github.com/mine/just-projecting/pkg/logger"
	msgbroker "github.com/mine/just-projecting/pkg/msg_broker"
	"github.com/mine/just-projecting/worker/handler"
)

type Worker interface {
	Run()
}

type KafkaWorker struct {
	kafkaConfig   msgbroker.KafkaConfig
	consumerGroup sarama.ConsumerGroup
	kafkaClient   sarama.Client
	api           *fiber.App
}

func NewKafkaWorker(ctx context.Context) Worker {
	log := logger.GetLoggerContext(ctx, "worker", "NewKafkaWorker")

	jsonByte, err := json.Marshal(config.Get("kafka"))
	if err != nil {
		log.Error(err)
		return nil
	}

	var cfg msgbroker.KafkaConfig
	err = json.Unmarshal(jsonByte, &cfg)
	if err != nil {
		log.Error(err)
		return nil
	}

	msgBroker, err := msgbroker.NewKafkaMsg(ctx)
	kafkaCli, consumerGroup := msgBroker.GetKafkaConn()
	return &KafkaWorker{
		kafkaConfig:   cfg,
		consumerGroup: consumerGroup,
		kafkaClient:   kafkaCli,
		api: fiber.New(
			fiber.Config{
				ErrorHandler: helper.Error,
			},
		),
	}
}

func (k *KafkaWorker) Run() {
	logrus.Info("--- Starting Kafka Consumer ---")
	ctx, cancel := context.WithCancel(context.Background())
	handler := handler.UserWorkerHandler{
		Ready: make(chan bool),
		Ctx:   ctx,
		Cfg:   k.kafkaConfig,
	}

	// HealthCheck Section
	gracefulStop := make(chan os.Signal)
	go k.healthCheck(gracefulStop)

	// Define topic List
	topics := []string{k.kafkaConfig.ExampleTopic}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := k.consumerGroup.Consume(ctx, topics, handler); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			if ctx.Err() != nil {
				return
			}

			handler.Ready = make(chan bool)
		}
	}()

	<-handler.Ready
	logrus.Info("Sarama Consumer up and Running")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}

	cancel()
	wg.Wait()
	if err := k.consumerGroup.Close(); err != nil {
		log.Panicf("Error closing consumer group: %v", err)
	}
	if err := k.kafkaClient.Close(); err != nil {
		log.Printf("Error closing kafka client: %v", err)
	}
}

func (k KafkaWorker) healthCheck(csig chan os.Signal) {
	k.api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	err := k.api.Listen(config.GetString("address") + ":" + config.GetString("worker_port"))
	if err != nil {
		log.Fatal("Failed start health check", err)
	}
	<-csig
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	k.api.Shutdown()
}
