/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package main

import (
	"fmt"
	"github.com/flanksource/github-app/cmd"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	root := &cobra.Command{
		Use:   "github-app",
		Short: "github-app : The flanksource github-app",
	}

	root.AddCommand(cmd.Serve,)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}



