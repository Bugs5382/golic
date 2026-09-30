package commands

/*
Copyright 2026 Shane & Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	_ "embed"

	"github.com/Bugs5382/golic/internal/logging"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var verbose bool

	// command
	var rootCmd = &cobra.Command{
		Use:   "golic",
		Short: "Inject, replace and remove license headers in source files",
		Long:  ``,
		// With no command, show the help and exit 0. Unknown commands are
		// rejected by cobra before this runs.
		RunE: func(cmd *cobra.Command, args []string) error {
			logging.L().Debug("no command given, printing help")
			return cmd.Help()
		},
	}

	// disable
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output (trace logging)")

	// sub commands
	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(InjectCmd())
	rootCmd.AddCommand(RemoveCmd())
	rootCmd.AddCommand(ReplaceCmd())

	return rootCmd
}
