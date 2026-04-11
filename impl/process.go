package impl

import (
	"context"
	"os"

	"github.com/AbsaOSS/golic/internal"
	"github.com/denormal/go-gitignore"
	"github.com/enescakir/emoji"
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

type Process struct {
	Ctx  context.Context
	Opts internal.Options

	cfg      *Config
	ignore   gitignore.GitIgnore
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
	log.Debugf("%s reading config path: %s", emoji.OpenBook, u.Opts.ConfigPath)
	log.Debugf("%s reading lic ignore path: %s", emoji.OpenBook, u.Opts.LicIgnore)
	log.Debugf("%s reading template: %s", emoji.OpenBook, u.Opts.Template)
	log.Debugf("%s reading search path: %s", emoji.OpenBook, u.Opts.SearchPath)

	if u.cfg, err = u.readCommonConfig(); err != nil {
		return
	}

	u.ignore, err = gitignore.NewFromFile(u.Opts.LicIgnore)
	if err != nil {
		return err
	}

	if _, err = os.Stat(u.Opts.ConfigPath); !os.IsNotExist(err) {
		log.Debugf("%s reading %s", emoji.OpenBook, aurora.BrightCyan(u.Opts.ConfigPath))
		log.Debugf("%s overriding %s with %s",
			emoji.ConstructionWorker, aurora.BrightCyan("master config"),
			aurora.BrightCyan(u.Opts.ConfigPath))
		if u.cfg, err = u.readLocalConfig(); err != nil {
			return
		}
	} else {
		if u.Opts.ConfigPath == "" {
			log.Debugf("%s no local found; using embeded.", emoji.FileFolder)
		} else {
			log.Debugf("%s skipping local %s", emoji.FileFolder, aurora.BrightCyan(u.Opts.ConfigPath))
		}
	}

	u.traverseFiles()

	return
}

func (u *Process) String() string {
	switch u.Opts.Type {
	case internal.LicenseInject:
		return aurora.BrightCyan("inject").String()
	case internal.LicenseRemove:
		return aurora.BrightCyan("remove").String()
	}
	return aurora.BrightRed("ERROR, unrecognised command").String()
}

func (u *Process) ExitCode() int {
	if u.Opts.ModifiedExitStatus && u.modified > 0 {
		return 1
	}
	return 0
}
