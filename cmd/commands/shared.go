package commands

/*
Apache License 2.0

Copyright 2026 Shane

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

	"github.com/Bugs5382/golic/cmd/logging"
	"github.com/Bugs5382/golic/internal"
	"github.com/spf13/cobra"
)

// setupAndValidate is the shared entry point for both commands
func setupAndValidate(cmd *cobra.Command, opts *internal.Options) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	logging.Init(verbose)

	configFlagChanged := cmd.Flags().Changed("config-path")

	if opts.ConfigPath != "" {
		_, err := os.Stat(opts.ConfigPath)
		if os.IsNotExist(err) {
			if configFlagChanged {
				return fmt.Errorf("custom config file not found: %s", opts.ConfigPath)
			}
			opts.ConfigPath = ""
		}
	}

	if opts.LicIgnore != "" {
		if _, err := os.Stat(opts.LicIgnore); os.IsNotExist(err) {
			return fmt.Errorf("custom ignore file not found: %s", opts.LicIgnore)
		}
	}

	if opts.Template == "" {
		return fmt.Errorf("license template not provided")
	}

	return nil
}
