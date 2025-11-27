package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	appConfigPath string
	vmConfigPath  string
	verbose       bool
)

var rootCmd = &cobra.Command{
	Use:   "vmdeploy",
	Short: "vmdeploy is a VM deployment engine",
	Long:  `vmdeploy is a Go-based deployment CLI for Linux VMs.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "[FATAL] %v\n", err)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&appConfigPath, "app-config", "appconfig.yaml", "Path to app config")
	rootCmd.PersistentFlags().StringVar(&vmConfigPath, "vm-config", "vmconfig.yaml", "Path to VM config")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(versionCmd)
}
