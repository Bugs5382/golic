package commands

/*
Apache License 2.0

Copyright 2026 Shane

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
	"fmt"

	"github.com/Bugs5382/golic"
	"github.com/Bugs5382/golic/impl"
	"github.com/Bugs5382/golic/internal"
	"github.com/spf13/cobra"
)

func RemoveCmd() *cobra.Command {

	opts := internal.Options{}

	opts.MasterConfig = golic.DefaultConfig

	// command
	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove licenses",
		Long:  ``,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return setupAndValidate(cmd, &opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// we are removing!
			opts.Type = 1

			// go ahead and start the inject process!
			i := impl.ProcessFile(cmd.Context(), opts)
			exitCode, err := internal.Command(i).MustRun()

			if err != nil {
				return err
			}

			if exitCode != 0 {
				return fmt.Errorf("something went wrong")
			}

			return nil
		},
	}

	// flags
	removeCmd.Flags().BoolVarP(&opts.ModifiedExitStatus, "modified-exit", "x", false,
		"If enabled, exits with status 1 when any file is modified. The settings is used by CI")
	removeCmd.Flags().BoolVarP(&opts.Dry, "dry", "d", false, "Dry run")

	removeCmd.Flags().StringVarP(&opts.Template, "template", "t", "", "License key")
	removeCmd.Flags().StringVarP(&opts.LicIgnore, "licignore", "l", ".licignore",
		".licignore path")
	removeCmd.Flags().StringVarP(&opts.Copyright, "copyright", "c",
		fmt.Sprintf("%d %s", internal.Year, "[Insert Company]"), "Copyright holder and year for the license header")
	removeCmd.Flags().StringVarP(&opts.ConfigPath, "config-path", "p", ".golic.yaml",
		"Path to the local configuration overriding config-url")
	//removeCmd.Flags().StringVarP(&internal.opts.SearchPath, "include-only", "i", "",
	//	"Used to execute only in reading into the path/directory provided")
	return removeCmd
}
