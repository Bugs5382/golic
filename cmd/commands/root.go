package commands

/*
Apache License 2.0

Copyright 2006 Shane

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
	"fmt"

	"github.com/Bugs5382/golic/impl"
	"github.com/spf13/cobra"
)

var BaseConfig *impl.Config

func RootCmd() *cobra.Command {
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

	// disable
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (e.g. tracing)")

	// sub commands
	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(InjectCmd())
	rootCmd.AddCommand(RemoveCmd())

	return rootCmd
}
