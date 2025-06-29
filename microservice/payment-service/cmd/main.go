package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "Multi-purpose CLI for app management",
}

func main() {
	rootCmd.AddCommand(restApi, etlConsumer)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
