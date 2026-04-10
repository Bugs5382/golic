package commands

import (
	"fmt"
	"os"

	"github.com/AbsaOSS/golic/helpers"
	"github.com/spf13/cobra"
)

func RemoveCmd() *cobra.Command {

	var removeOptions helpers.Options

	// command
	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove licenses",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// golic config
			configPath := helpers.RemoveOptions.ConfigPath
			if configPath == "" {
				configPath = ".golic.yaml"
			}

			// Check if the file actually exists
			if _, err := os.Stat(configPath); err != nil {
				if os.IsNotExist(err) {
					_ = cmd.Help() // Print usage instructions
					return fmt.Errorf("config file not found: ensure '.golic.yaml' exists in the current" +
						" directory or provide a valid path")
				}
				// Catch any other potential file system errors (e.g., permission denied)
				return fmt.Errorf("error accessing config file: %w", err)
			}

			// Ensure the resolved path is saved back to your options so downstream code uses it
			removeOptions.ConfigPath = configPath

			// ignore lic
			ignorePath := helpers.RemoveOptions.LicIgnore
			if ignorePath == "" {
				ignorePath = ".licignore"
			}

			// Check if the file actually exists
			if _, err := os.Stat(ignorePath); err != nil {
				if os.IsNotExist(err) {
					_ = cmd.Help() // Print usage instructions
					return fmt.Errorf("ignore file not found: ensure '.licignore' exists in the current" +
						" directory or provide a valid path")
				}
				// Catch any other potential file system errors (e.g., permission denied)
				return fmt.Errorf("error accessing config file: %w", err)
			}

			// Ensure the resolved path is saved back to your options so downstream code uses it
			removeOptions.LicIgnore = ignorePath

			templateSelected := helpers.RemoveOptions.Template
			print(templateSelected)
			if templateSelected == "" {
				_ = cmd.Help() // Print usage instructions
				return fmt.Errorf("licence template not provided")
			}

			// dry run options
			removeOptions.Dry = helpers.RemoveOptions.Dry
			// modified status
			removeOptions.ModifiedExitStatus = helpers.RemoveOptions.ModifiedExitStatus
			// search dir settings
			removeOptions.SearchPath = helpers.RemoveOptions.SearchPath
			// we are removing!
			removeOptions.Type = 1

			// go ahead and start the remove process!

			//	fmt.Printf(" %s %s\n", emoji.Rocket, aurora.BrightWhite("done"))
			//	fmt.Printf(" %s %s\n", emoji.FaceScreamingInFear, aurora.BrightWhite("found files with missing a license, exit"))

			return nil
		},
	}

	// flags
	removeCmd.Flags().BoolVarP(&helpers.RemoveOptions.Dry, "dry", "d", false, "dry run")
	removeCmd.Flags().StringVarP(&helpers.RemoveOptions.LicIgnore, "licignore", "l", ".licignore", ".licignore path")
	removeCmd.Flags().StringVarP(&helpers.RemoveOptions.Template, "template", "t", "", "license key")
	removeCmd.Flags().StringVarP(&helpers.RemoveOptions.Copyright, "copyright", "c",
		fmt.Sprintf("%d %s", helpers.Year, helpers.Company), "copyright holder and year for the license header",
	)
	removeCmd.Flags().StringVarP(&helpers.RemoveOptions.ConfigPath, "config-path", "p", ".golic.yaml", "path to the local configuration overriding config-url")

	return removeCmd
}
