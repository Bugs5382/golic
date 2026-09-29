package impl

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
	"context"
	"os"

	"github.com/Bugs5382/golic/internal"
	"github.com/enescakir/emoji"
	"github.com/goccy/go-yaml"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog/log"
)

type Process struct {
	Ctx  context.Context
	Opts internal.Options

	cfgBase  *Config
	cfg      *Config
	scope    *scope
	modified int
}

func ProcessFile(ctx context.Context, options internal.Options) *Process {
	return &Process{
		Ctx:  ctx,
		Opts: options,

		modified: 0,
	}
}

func (u *Process) Run() (err error) {
	// debug commands
	log.Debug().Msgf("%s reading config path: %s", emoji.OpenBook, u.Opts.ConfigPath)
	log.Debug().Msgf("%s reading lic ignore path: %s", emoji.OpenBook, u.Opts.LicIgnore)
	log.Debug().Msgf("%s reading template: %s", emoji.OpenBook, u.Opts.Template)
	log.Debug().Msgf("%s reading search path: %s", emoji.OpenBook, u.Opts.SearchPath)

	if u.cfgBase, err = u.readCommonConfig(); err != nil {
		return
	}

	if _, err = os.Stat(u.Opts.ConfigPath); !os.IsNotExist(err) {
		log.Debug().Msgf("%s reading %s", emoji.OpenBook, aurora.BrightCyan(u.Opts.ConfigPath))
		log.Debug().Msgf("%s merging %s with %s", emoji.ConstructionWorker, aurora.BrightCyan(u.Opts.ConfigPath), aurora.BrightCyan("master config"))
		if u.cfg, err = u.readLocalConfig(); err != nil {
			return
		}
	} else {
		if u.Opts.ConfigPath == "" {
			log.Debug().Msgf("%s no local found; using embeded.", emoji.FileFolder)
		} else {
			log.Debug().Msgf("%s skipping local %s", emoji.FileFolder, aurora.BrightCyan(u.Opts.ConfigPath))
		}
		u.cfg = u.cfgBase
	}

	if log.Debug().Enabled() {
		// Marshal the merged config to YAML for a "pretty-print" effect
		confBytes, _ := yaml.Marshal(u.cfg)

		log.Debug().
			Msgf("Final Configuration Loaded:\n---\n%s\n---", string(confBytes))
	}

	log.Debug().Strs("ignore", u.cfg.Golic.Ignore).Msg("config ignore patterns in effect")
	if u.scope, err = newScope(u.Opts.LicIgnore, u.cfg.Golic.Ignore); err != nil {
		return err
	}

	err = u.traverseFiles()

	return
}

func (u *Process) String() string {
	switch u.Opts.Type {
	case internal.LicenseInject:
		return aurora.BrightCyan("inject").String()
	case internal.LicenseRemove:
		return aurora.BrightCyan("remove").String()
	case internal.LicenseReplace:
		return aurora.BrightCyan("replace").String()
	}
	return aurora.BrightRed("ERROR, unrecognised command").String()
}

// Changes returns a *internal.ChangesError when -x is set and at least one
// file was modified (or, in a dry run, would be). It returns nil otherwise.
func (u *Process) Changes() error {
	if !u.Opts.ModifiedExitStatus {
		log.Debug().Int("modified", u.modified).Msg("modified-exit not set; changes do not affect the exit status")
		return nil
	}
	if u.modified == 0 {
		log.Debug().Msg("modified-exit set and no file needs a change")
		return nil
	}
	log.Debug().Int("modified", u.modified).Bool("dry", u.Opts.Dry).Msg("modified-exit set and files need a change")
	return &internal.ChangesError{Type: u.Opts.Type, Dry: u.Opts.Dry, Count: u.modified}
}
