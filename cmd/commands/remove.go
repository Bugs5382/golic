package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/AbsaOSS/golic/internal"
	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
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

			verbose, _ := cmd.Flags().GetBool("verbose")

			if verbose {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}

			// golic config
			configPath := internal.RemoveOptions.ConfigPath
			if configPath == "" {
				configPath = ".golic.yaml"
			}

			// Check if the file actually exists
			if _, err := os.Stat(configPath); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("config file not found: ensure '.golic.yaml' exists in the current" +
						" directory or provide a valid path")
				}
				// Catch any other potential file system errors (e.g., permission denied)
				return fmt.Errorf("error accessing config file: %w", err)
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

			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Suffix = "Removing licenses, please wait..."
			s.Start()

			// go ahead and start the remove process!
			//i := impl.NewRemove(cmd.Context(), removeOptions)
			//exitCode := pkg.Command(i).MustRun()
			//
			//s.Stop()
			//
			//if exitCode != 0 {
			//	// Handle error
			//}
			return nil
		},
	}

	// flags
	removeCmd.Flags().BoolVarP(&internal.RemoveOptions.Dry, "dry", "d", false, "dry run")
	removeCmd.Flags().StringVarP(&internal.RemoveOptions.LicIgnore, "licignore", "l", ".licignore", ".licignore path")
	removeCmd.Flags().StringVarP(&internal.RemoveOptions.Template, "template", "t", "", "license key")
	removeCmd.Flags().StringVarP(&internal.RemoveOptions.Copyright, "copyright", "c",
		fmt.Sprintf("%d %s", internal.Year, internal.Company), "copyright holder and year for the license header",
	)
	removeCmd.Flags().StringVarP(&internal.RemoveOptions.ConfigPath, "config-path", "p", ".golic.yaml", "path to the local configuration overriding config-url")

	return removeCmd
}
