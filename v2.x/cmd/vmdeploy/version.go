package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CLIVersion is the reported version of the vmdeploy CLI.
// For now, it's static. Later you can inject it at build time via -ldflags.
const CLIVersion = "v2.0.5"

// versionCmd implements: vmdeploy version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the vmdeploy CLI version",
	Long: `Shows the vmdeploy version.

In future versions this will expose:
  - CLI version
  - Go runtime version
  - Git commit
  - Build date
  - Platform`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vmdeploy CLI version: ", CLIVersion)
	},
}
