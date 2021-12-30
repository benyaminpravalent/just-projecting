package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/mine/just-projecting/api"
	"github.com/mine/just-projecting/pkg/config"
	"github.com/mine/just-projecting/pkg/database"
	"github.com/mine/just-projecting/pkg/logger"
)

var configURL string
var serverCommand = &cobra.Command{
	Use: "serve",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
	PostRun: func(cmd *cobra.Command, args []string) {

	},
}

func startServer() {
	ctx := context.Background()

	if err := config.Load(DefaultConfig, configURL); err != nil {
		log.Fatal(err)
	}

	logger.Configure()

	database.InitMySql(ctx)

	srv := api.NewFiberServer()
	if err := srv.Configure(ctx); err != nil {
		log.Fatal(err)
	}
	if err := srv.Serve(); err != nil {
		log.Fatal(err)
	}
}
