package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "just-projecting",
	Short: "Just Projecting",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use serve to start a server")
		fmt.Println("Use -h to see the list of command")
	},
}

func Run() {
	serverCommand.PersistentFlags().StringVarP(&configURL, "config", "c", "", "Config URL i.e. file://config.json")
	kafkaConsumerCmd.PersistentFlags().StringVarP(&configURL, "config", "c", "", "Config URL i.e. file://config.json")
	rootCommand.AddCommand(serverCommand)
	rootCommand.AddCommand(kafkaConsumerCmd)

	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
