package commands

import (
	"fmt"

	"github.com/AbsaOSS/golic/helpers"
	"github.com/spf13/cobra"
)

func InjectCmd() *cobra.Command {

	// command
	var injectCmd = &cobra.Command{
		Use:   "inject",
		Short: "Injects licenses",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			//if _, err := os.Stat(injectOptions.LicIgnore); os.IsNotExist(err) {
			//	//logger.Error().Msgf("invalid license path '%s'", injectOptions.LicIgnore)
			//	_ = cmd.Help()
			//	os.Exit(1)
			//}
			//if masterconfig == "" {
			//	logger.Error().Msgf("invalid master config. ")
			//	_ = cmd.Help()
			//	os.Exit(1)
			//}
			//injectOptions.MasterConfig = masterconfig
			//injectOptions.Type = update.LicenseInject
			//i := update.New(ctx, injectOptions)
			//exitCode = temp.Command(i).MustRun()
			//if exitCode == 0 {
			//	fmt.Printf(" %s %s\n", emoji.Rocket, aurora.BrightWhite("done"))
			//} else {
			//	fmt.Printf(" %s %s\n", emoji.FaceScreamingInFear, aurora.BrightWhite("found files with missing a license, exit"))
			//}
		},
	}

	// flags
	injectCmd.Flags().BoolVarP(&helpers.InjectOptions.ModifiedExitStatus, "modified-exit", "x", false,
		"If enabled, exits with status 1 when any file is modified. The settings is used by CI")
	injectCmd.Flags().BoolVarP(&helpers.InjectOptions.Dry, "dry", "d", false, "dry run")

	injectCmd.Flags().StringVarP(&helpers.InjectOptions.Template, "template", "t", "", "license key")
	injectCmd.Flags().StringVarP(&helpers.InjectOptions.LicIgnore, "licignore", "l", ".licignore",
		".licignore path")
	injectCmd.Flags().StringVarP(
		&helpers.InjectOptions.Copyright,
		"copyright",
		"c",
		fmt.Sprintf("%d %s", helpers.Year, helpers.Company),
		"copyright holder and year for the license header",
	)
	injectCmd.Flags().StringVarP(&helpers.InjectOptions.ConfigPath, "config-path", "p", ".golic.yaml",
		"path to the local configuration overriding config-url")

	return injectCmd
}
