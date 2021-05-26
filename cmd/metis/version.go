package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"oss.navercorp.com/metis/metis-server/internal/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Metis Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Metis Server: %s\n", version.Version)
			fmt.Printf("Commit: %s\n", version.GitCommit)
			fmt.Printf("Go: %s\n", runtime.Version())
			fmt.Printf("Build date: %s\n", version.BuildDate)
			return nil
		},
	}
}

func init() {
	rootCmd.AddCommand(newVersionCmd())
}
