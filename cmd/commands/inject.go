package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/AbsaOSS/golic/impl"
	"github.com/AbsaOSS/golic/internal"
	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func InjectCmd(masterConfig string) *cobra.Command {

	var injectOptions internal.Options
	injectOptions.MasterConfig = masterConfig

	// command
	var injectCmd = &cobra.Command{
		Use:   "inject",
		Short: "Injects licenses",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {

			verbose, _ := cmd.Flags().GetBool("verbose")

			if verbose {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}

			// golic config
			configPath := internal.InjectOptions.ConfigPath
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
			injectOptions.ConfigPath = configPath

			// ignore lic
			ignorePath := internal.InjectOptions.LicIgnore
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
			injectOptions.LicIgnore = ignorePath

			templateSelected := internal.InjectOptions.Template
			if templateSelected == "" {
				return fmt.Errorf("licence template not provided")
			}

			// template setting
			injectOptions.Template = templateSelected
			// dry run options
			injectOptions.Dry = internal.InjectOptions.Dry
			// modified status
			injectOptions.ModifiedExitStatus = internal.InjectOptions.ModifiedExitStatus
			// search dir settings
			injectOptions.SearchPath = internal.InjectOptions.SearchPath
			// we are injecting!
			injectOptions.Type = 0
			// verbose
			injectOptions.Verbose = verbose

			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Suffix = "Injecting licenses, please wait..."
			s.Start()

			// go ahead and start the inject process!
			i := impl.NewInject(cmd.Context(), injectOptions)
			exitCode := internal.Command(i).MustRun()

			s.Stop()

			if exitCode != 0 {
				return fmt.Errorf("something went wrong")
			}

			return nil
		},
	}

	// flags
	injectCmd.Flags().BoolVarP(&internal.InjectOptions.ModifiedExitStatus, "modified-exit", "x", false,
		"If enabled, exits with status 1 when any file is modified. The settings is used by CI")
	injectCmd.Flags().BoolVarP(&internal.InjectOptions.Dry, "dry", "d", false, "Dry run")

	injectCmd.Flags().StringVarP(&internal.InjectOptions.Template, "template", "t", "", "License key")
	injectCmd.Flags().StringVarP(&internal.InjectOptions.LicIgnore, "licignore", "l", ".licignore",
		".licignore path")
	injectCmd.Flags().StringVarP(&internal.InjectOptions.Copyright, "copyright", "c",
		fmt.Sprintf("%d %s", internal.Year, internal.Company), "Copyright holder and year for the license header")
	injectCmd.Flags().StringVarP(&internal.InjectOptions.ConfigPath, "config-path", "p", ".golic.yaml",
		"Path to the local configuration overriding config-url")
	injectCmd.Flags().StringVarP(&internal.InjectOptions.SearchPath, "include-only", "i", "",
		"Used to execute only in reading into the path/directory provided")

	return injectCmd
}
