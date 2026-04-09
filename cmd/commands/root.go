package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var verbose bool

	// command
	var rootCmd = &cobra.Command{
		Use:           "golic",
		Short:         "golic license injector",
		Long:          ``,
		SilenceUsage:  true, // Prevents automatic help print on error
		SilenceErrors: true, // We will handle the error printing in main()
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// cmd.Printf("Executing: %s %v\n", cmd.CommandPath(), args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				_ = cmd.Help()
				return fmt.Errorf("%s", "no arguments passed")
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
	}

	// flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (e.g. tracing)")

	// sub commands
	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(InjectCmd())
	rootCmd.AddCommand(RemoveCmd())

	return rootCmd
}
