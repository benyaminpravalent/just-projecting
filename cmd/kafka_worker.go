package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/mine/just-projecting/pkg/config"
	"github.com/mine/just-projecting/pkg/database"
	"github.com/mine/just-projecting/pkg/logger"
	"github.com/mine/just-projecting/worker"
)

// kafkaConsumerCmd represents the kafkaConsumer command
var kafkaConsumerCmd = &cobra.Command{
	Use:   "kafka-worker",
	Short: "Run kafka worker",
	Run: func(cmd *cobra.Command, args []string) {
		startWorker()
	},
}

func startWorker() {
	ctx := context.Background()
	if err := config.Load(DefaultConfig, configURL); err != nil {
		log.Fatal(err)
	}

	logger.Configure()

	database.InitMySql(ctx)

	wk := worker.NewKafkaWorker(ctx)
	wk.Run()
}
