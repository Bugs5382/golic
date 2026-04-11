package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RootCmd(masterConfig string) *cobra.Command {
	var verbose bool

	// command
	var rootCmd = &cobra.Command{
		Use:   "golic",
		Short: "golic license injector",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("%s", "no arguments passed")
			}
			return nil
		},
	}

	// flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (e.g. tracing)")

	// sub commands
	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(InjectCmd(masterConfig))
	rootCmd.AddCommand(RemoveCmd(masterConfig))

	return rootCmd
}
