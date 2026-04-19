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

	f.BoolVarP(&opts.ModifiedExitStatus, "modified-exit", "x", false, "If enabled, exits with status 1 when any file is modified. The settings is used by CI")
	f.BoolVarP(&opts.Dry, "dry", "d", false, "Dry run")

	f.StringVarP(&opts.Template, "template", "t", "", "License key")
	f.StringVarP(&opts.LicIgnore, "licignore", "l", ".licignore", ".licignore path")
	f.StringVarP(&opts.Copyright, "copyright", "c", fmt.Sprintf("%d %s", internal.Year, "[Insert Company]"), "Copyright holder and year for the license header")
	f.StringVarP(&opts.ConfigPath, "config-path", "p", ".golic.yaml", "Path to the local configuration overriding config-url")
}
