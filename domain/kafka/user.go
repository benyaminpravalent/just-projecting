package kafka

import (
	"context"
	"encoding/json"

	"github.com/mine/just-projecting/domain/model"
	"github.com/mine/just-projecting/pkg/config"
	"github.com/mine/just-projecting/pkg/logger"
	msgbroker "github.com/mine/just-projecting/pkg/msg_broker"
)

type UserKafkaRepository interface {
	Publish(ctx context.Context, user model.User) error
}

type userKafka struct {
	msgBroker   msgbroker.MsgBroker
	kafkaConfig msgbroker.KafkaConfig
}

func NewUserKafka(ctx context.Context) UserKafkaRepository {
	log := logger.GetLoggerContext(ctx, "kafkaDomain", "NewUserKafka")
	msgBroker, err := msgbroker.NewKafkaMsg(ctx)
	if err != nil {
		log.Errorf("Failed getting new kafka function : %s", err.Error())
		return nil
	}

	jsonByte, err := json.Marshal(config.Get("kafka"))
	if err != nil {
		log.Error(err)
		return nil
	}

	var kfkCfg msgbroker.KafkaConfig
	err = json.Unmarshal(jsonByte, &kfkCfg)

	return &userKafka{
		msgBroker:   msgBroker,
		kafkaConfig: kfkCfg,
	}
}

func (k *userKafka) Publish(ctx context.Context, user model.User) error {
	log := logger.GetLoggerContext(ctx, "msgBroker", "NewKafka")
	by, err := json.Marshal(user)
	if err != nil {
		log.Error(err)
		return err
	}

	return k.msgBroker.Publish(ctx, k.kafkaConfig.ExampleTopic, string(by))
}
