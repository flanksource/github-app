/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package cmd

import (
	"fmt"
	"github.com/flanksource/github-app/config"
	"github.com/flanksource/github-app/server"
	"github.com/spf13/cobra"
)



var (
	Runners = &cobra.Command{
		Use: "runners",
		Short: "commands related to managing github runners",
	}
	Cleanup = &cobra.Command{
	Use:   "cleanup",
	Short: "cleans up offline runners",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.ReadConfig(configFile)
		if err != nil {
			return fmt.Errorf("error reading config file %v: %v", configFile, err)
		}

		server, err := server.New(config)
		if err != nil {
			return fmt.Errorf("error starting server: %v", err)
		}

		return server.Start()
	},
}
)

func init() {
	Runners.Flags().StringVar(&configFile, "configuration file", "config.yaml", "The config file containing secrets, endpoints, etc.")
}
