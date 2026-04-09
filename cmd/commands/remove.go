package commands

import (
	"fmt"

	"github.com/AbsaOSS/golic/helpers"
	"github.com/spf13/cobra"
)

func RemoveCmd() *cobra.Command {
	// command
	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove licenses",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			//if _, err := os.Stat(commands.injectOptions.LicIgnore); os.IsNotExist(err) {
			//	logger.Error().Msgf("invalid license path '%s'", commands.injectOptions.LicIgnore)
			//	_ = cmd.Help()
			//	os.Exit(1)
			//}
			//commands.injectOptions.MasterConfig = masterconfig
			//commands.injectOptions.Type = update.LicenseRemove
			//i := update.New(ctx, commands.injectOptions)
			//_ = Command(i).MustRun()
			//fmt.Printf(" %s %s\n", emoji.Rocket, aurora.BrightWhite("done"))
		},
	}

	// flags
	removeCmd.Flags().BoolVarP(&helpers.InjectOptions.Dry, "dry", "d", false, "dry run")
	removeCmd.Flags().StringVarP(&helpers.InjectOptions.LicIgnore, "licignore", "l", ".licignore", ".licignore path")
	removeCmd.Flags().StringVarP(&helpers.InjectOptions.Template, "template", "t", "apache2", "license key")
	removeCmd.Flags().StringVarP(&helpers.InjectOptions.Copyright, "copyright",
		"c",
		fmt.Sprintf("%d %s", helpers.Year, helpers.Company),
		"copyright holder and year for the license header",
	)
	removeCmd.Flags().StringVarP(&helpers.InjectOptions.ConfigPath, "config-path", "p", ".golic.yaml", "path to the local configuration overriding config-url")

	return removeCmd
}
