package impl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AbsaOSS/golic/internal"
	"github.com/denormal/go-gitignore"
	"github.com/enescakir/emoji"
	"github.com/goccy/go-yaml"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog/log"
)

// readCommonConfig Read the commong/master config
func (u *Process) readCommonConfig() (c *Config, err error) {
	c = &Config{}
	err = yaml.Unmarshal([]byte(u.Opts.MasterConfig), c)
	return
}

// readLocalConfig Read the local config.
func (u *Process) readLocalConfig() (*Config, error) {
	var c = &Config{}
	var rc = *u.cfg
	yamlFile, err := os.ReadFile(u.Opts.ConfigPath)
	if err != nil {
		return nil, nil
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, nil
	}
	for k, v := range c.Golic.Licenses {
		rc.Golic.Licenses[k] = v
	}
	for k, v := range c.Golic.Rules {
		rc.Golic.Rules[k] = v
	}
	return &rc, nil
}

// traverseFiles Go through all files in paths and process. Will ignore files and folders that match GitIgnore.
func (u *Process) traverseFiles() {
	skipped := 0
	visited := 0
	p := func(path string, i gitignore.GitIgnore, o internal.Options, config *Config) (err error) {
		if !i.Ignore(path) {
			var skip bool
			symbol := ""
			cp := aurora.BrightYellow(path)
			visited++
			if err, skip = processUpdate(path, o, config); skip {
				symbol = "-> skip"
				cp = aurora.Magenta(path)
				skipped++
			}
			_, _ = emoji.Printf(" %s  %s %s  \n", emoji.Minus, cp, aurora.BrightMagenta(symbol))
		}
		return
	}

	err := filepath.Walk("./",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return p(path, u.ignore, u.Opts, u.cfg)
			}
			return nil
		})

	if err != nil {
		log.Error()
	}

	u.modified = visited - skipped
	displaySummary(skipped, visited)
}

// processUpdate Update the file, but how?
func processUpdate(path string, o internal.Options, config *Config) (err error, skip bool) {
	switch o.Type {
	case internal.LicenseInject:
		//return injectFile(path, o, config)
	case internal.LicenseRemove:
		//return removeFile(path, o, config)
	}
	return fmt.Errorf("invalid license type"), true
}

func displaySummary(skipped, visited int) {
	if skipped == visited {
		fmt.Printf("\n %s %v/%v %s\n\n", emoji.Ice, aurora.BrightCyan(visited-skipped), aurora.BrightWhite(visited), aurora.BrightCyan("changed"))
		return
	}
	fmt.Printf("\n %s %v/%v %s\n\n", emoji.Fire, aurora.BrightYellow(visited-skipped), aurora.BrightWhite(visited), aurora.BrightYellow("changed"))
}
