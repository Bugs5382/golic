package commands

/*
Apache License 2.0

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
	"fmt"

	"github.com/Bugs5382/golic/internal"
	"github.com/spf13/cobra"
)

func addCommonFlags(cmd *cobra.Command, opts *internal.Options) {
	f := cmd.Flags()

	f.BoolVarP(&opts.ModifiedExitStatus, "modified-exit", "x", false, "Exit with status 1 when any file is modified, or would be in a dry run (for CI). Errors exit with 2")
	f.BoolVarP(&opts.Dry, "dry", "d", false, "Report what would change without writing any files")

	f.StringVarP(&opts.Template, "template", "t", "", "License key from the config to use (for example apache2 or mit)")
	f.StringVarP(&opts.LicIgnore, "licignore", "l", ".licignore", "Path to the .licignore file that selects which files to process")
	f.StringVarP(&opts.Copyright, "copyright", "c", fmt.Sprintf("%d %s", internal.Year, "[Insert Company]"), "Copyright holder and year for the license header")
	f.StringVarP(&opts.ConfigPath, "config-path", "p", ".golic.yaml", "Path to a local config merged over the built-in rules and licenses")
}

// addLicenseFileFlag adds --license-file to the commands that write headers.
// With no value it names LICENSE in the directory golic runs in.
func addLicenseFileFlag(cmd *cobra.Command, opts *internal.Options) {
	f := cmd.Flags()
	f.StringVar(&opts.LicenseFile, "license-file", "", "Also write the full official licence text for the template to this file (LICENSE when no path is given), and count a missing or different file as a change")
	f.Lookup("license-file").NoOptDefVal = "LICENSE"
}
