package commands

/*
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
	"github.com/Bugs5382/golic/internal"
	"github.com/spf13/cobra"
)

func ReplaceCmd() *cobra.Command {

	opts := internal.Options{}

	// command
	var replaceCmd = &cobra.Command{
		Use:   "replace",
		Short: "Replace the existing license header with the configured one",
		Long:  `Swap the license header golic wrote for the configured one in place, in a single pass. Files with no header are stamped as inject would; other license headers are left alone.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return setupAndValidate(cmd, &opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Type = internal.LicenseReplace
			return runProcess(cmd, opts)
		},
	}

	// flags
	addCommonFlags(replaceCmd, &opts)
	addLicenseFileFlag(replaceCmd, &opts)
	return replaceCmd
}
