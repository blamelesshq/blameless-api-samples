package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "Command line interface for SLO API example",
		Long:  `Command line interface for SLO API example`,
	}

	rootCmd.AddCommand(sli())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("unable to start command line \n%+v", err)
	}
}
