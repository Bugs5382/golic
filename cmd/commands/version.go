package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.8.0"

// VersionCmd Show the Golic Version
func VersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Golic",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
	return versionCmd
}
