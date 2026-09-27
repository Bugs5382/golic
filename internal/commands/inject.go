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
	"github.com/Bugs5382/golic/internal"
	"github.com/spf13/cobra"
)

func InjectCmd() *cobra.Command {

	opts := internal.Options{}

	// command
	var injectCmd = &cobra.Command{
		Use:   "inject",
		Short: "Add the license header to files that are missing it",
		Long:  ``,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return setupAndValidate(cmd, &opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Type = internal.LicenseInject
			return runProcess(cmd, opts)
		},
	}

	// flags
	addCommonFlags(injectCmd, &opts)
	return injectCmd
}
