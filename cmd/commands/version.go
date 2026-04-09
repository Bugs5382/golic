package commands

import (
	"fmt"

	"github.com/AbsaOSS/golic/cmd/logging"
	"github.com/spf13/cobra"
)

const version = "v0.8.0"

// VersionCmd Show the Golic Version
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Golic",
	Long:  "",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging.LogCommandExecution(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
