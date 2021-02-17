package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "metis [options]",
}

func Run() int {
	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}