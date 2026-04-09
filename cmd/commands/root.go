package commands

import (
	"fmt"
	"os"

	"github.com/AbsaOSS/golic/cmd/logging"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var verbose bool

	// command
	var rootCmd = &cobra.Command{
		Use:   "golic",
		Short: "golic license injector",
		Long:  ``,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logging.LogCommandExecution(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				_ = cmd.Help()
				return fmt.Errorf("no parameters included")
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	// flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (e.g. tracing)")

	// sub commands
	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(InjectCmd())
	rootCmd.AddCommand(RemoveCmd())

	return rootCmd
}
