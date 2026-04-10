package commands

import (
	"fmt"
	"os"

	"github.com/AbsaOSS/golic/helpers"
	"github.com/spf13/cobra"
)

func InjectCmd() *cobra.Command {

	var injectOptions helpers.Options

	// command
	var injectCmd = &cobra.Command{
		Use:   "inject",
		Short: "Injects licenses",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {

			// golic config
			configPath := helpers.InjectOptions.ConfigPath
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
			injectOptions.ConfigPath = configPath

			// ignore lic
			ignorePath := helpers.InjectOptions.LicIgnore
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
			injectOptions.LicIgnore = ignorePath

			templateSelected := helpers.InjectOptions.Template
			print(templateSelected)
			if templateSelected == "" {
				_ = cmd.Help() // Print usage instructions
				return fmt.Errorf("licence template not provided")
			}

			// dry run options
			injectOptions.Dry = helpers.InjectOptions.Dry
			// modified status
			injectOptions.ModifiedExitStatus = helpers.InjectOptions.ModifiedExitStatus
			// search dir settings
			injectOptions.SearchPath = helpers.InjectOptions.SearchPath
			// we are injecting!
			injectOptions.Type = 0

			// go ahead and start the inject process!

			//	fmt.Printf(" %s %s\n", emoji.Rocket, aurora.BrightWhite("done"))
			//	fmt.Printf(" %s %s\n", emoji.FaceScreamingInFear, aurora.BrightWhite("found files with missing a license, exit"))

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	// flags
	injectCmd.Flags().BoolVarP(&helpers.InjectOptions.ModifiedExitStatus, "modified-exit", "x", false,
		"If enabled, exits with status 1 when any file is modified. The settings is used by CI")
	injectCmd.Flags().BoolVarP(&helpers.InjectOptions.Dry, "dry", "d", false, "Dry run")

	injectCmd.Flags().StringVarP(&helpers.InjectOptions.Template, "template", "t", "", "License key")
	injectCmd.Flags().StringVarP(&helpers.InjectOptions.LicIgnore, "licignore", "l", ".licignore",
		".licignore path")
	injectCmd.Flags().StringVarP(&helpers.InjectOptions.Copyright, "copyright", "c",
		fmt.Sprintf("%d %s", helpers.Year, helpers.Company), "Copyright holder and year for the license header")
	injectCmd.Flags().StringVarP(&helpers.InjectOptions.ConfigPath, "config-path", "p", ".golic.yaml",
		"Path to the local configuration overriding config-url")
	injectCmd.Flags().StringVarP(&helpers.InjectOptions.SearchPath, "include-only", "i", "",
		"Used to execute only in reading into the path/directory provided")

	return injectCmd
}
