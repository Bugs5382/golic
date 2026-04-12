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
	"fmt"
	"os"

	"github.com/AbsaOSS/golic/impl"
	"github.com/AbsaOSS/golic/internal"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func RemoveCmd(masterConfig string) *cobra.Command {

	var removeOptions internal.Options
	removeOptions.MasterConfig = masterConfig

	// command
	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove licenses",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {

			levelStr := os.Getenv("LOG_LEVEL")
			verbose, _ := cmd.Flags().GetBool("verbose")
			if levelStr == "" {
				if verbose {
					zerolog.SetGlobalLevel(zerolog.DebugLevel)
				} else {
					zerolog.SetGlobalLevel(zerolog.InfoLevel)
				}
			}

			// golic config
			configPath := internal.InjectOptions.ConfigPath
			if configPath != "" {
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					return fmt.Errorf("config file not found: ensure '.golic.yaml' exists or check path: %s", configPath)
				} else if err != nil {
					return fmt.Errorf("error accessing config file: %w", err)
				}
			}

			// Ensure the resolved path is saved back to your options so downstream code uses it
			removeOptions.ConfigPath = configPath

			// ignore lic
			ignorePath := internal.RemoveOptions.LicIgnore
			if ignorePath == "" {
				ignorePath = ".licignore"
			}

			// Check if the file actually exists
			if _, err := os.Stat(ignorePath); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("ignore file not found: ensure '.licignore' exists in the current" +
						" directory or provide a valid path")
				}
				// Catch any other potential file system errors (e.g., permission denied)
				return fmt.Errorf("error accessing config file: %w", err)
			}

			// Ensure the resolved path is saved back to your options so downstream code uses it
			removeOptions.LicIgnore = ignorePath

			templateSelected := internal.RemoveOptions.Template
			if templateSelected == "" {
				return fmt.Errorf("licence template not provided")
			}

			// copyright string
			removeOptions.Copyright = internal.RemoveOptions.Copyright
			// template setting
			removeOptions.Template = templateSelected
			// dry run options
			removeOptions.Dry = internal.RemoveOptions.Dry
			// modified status
			removeOptions.ModifiedExitStatus = internal.RemoveOptions.ModifiedExitStatus
			// search dir settings
			removeOptions.SearchPath = internal.RemoveOptions.SearchPath
			// we are removing!
			removeOptions.Type = 1
			// verbose
			removeOptions.Verbose = verbose

			// go ahead and start the inject process!
			i := impl.ProcessFile(cmd.Context(), removeOptions)
			exitCode := internal.Command(i).MustRun()

			if exitCode != 0 {
				return fmt.Errorf("something went wrong")
			}

			return nil
		},
	}

	// flags
	removeCmd.Flags().BoolVarP(&internal.RemoveOptions.ModifiedExitStatus, "modified-exit", "x", false,
		"If enabled, exits with status 1 when any file is modified. The settings is used by CI")
	removeCmd.Flags().BoolVarP(&internal.RemoveOptions.Dry, "dry", "d", false, "Dry run")

	removeCmd.Flags().StringVarP(&internal.RemoveOptions.Template, "template", "t", "", "License key")
	removeCmd.Flags().StringVarP(&internal.RemoveOptions.LicIgnore, "licignore", "l", ".licignore",
		".licignore path")
	removeCmd.Flags().StringVarP(&internal.RemoveOptions.Copyright, "copyright", "c",
		fmt.Sprintf("%d %s", internal.Year, "[Insert Company]"), "Copyright holder and year for the license header")
	removeCmd.Flags().StringVarP(&internal.RemoveOptions.ConfigPath, "config-path", "p", "",
		"Path to the local configuration overriding config-url")
	removeCmd.Flags().StringVarP(&internal.RemoveOptions.SearchPath, "include-only", "i", "",
		"Used to execute only in reading into the path/directory provided")
	return removeCmd
}
