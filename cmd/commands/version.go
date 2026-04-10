package commands

import (
	"github.com/AbsaOSS/golic/helpers"
	"github.com/spf13/cobra"
)

// VersionCmd Show the Golic Version
func VersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Golic",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(helpers.Version)
		},
	}
	return versionCmd
}
