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

var configFile string

var Serve = &cobra.Command{
	Use:   "serve",
	Short: "starts a github-app server",
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

func init() {
	Serve.Flags().StringVar(&configFile, "configuration file", "config.yaml", "The config file containing secrets, endpoints, etc.")
}
