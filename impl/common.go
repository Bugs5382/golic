package impl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AbsaOSS/golic/internal"
	"github.com/AbsaOSS/golic/pkg"
	"github.com/denormal/go-gitignore"
	"github.com/enescakir/emoji"
	"github.com/goccy/go-yaml"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog/log"
)

func injectFile(path string, o internal.Options, config *Config) (err error, skip bool) {
	source, err := read(path)
	if err != nil {
		return err, false
	}
	rule := getRule(config, path)
	license, err := getCommentedLicense(config, o, rule)
	if err != nil {
		return err, false
	}
	// license is injected, continue
	if strings.Contains(source, license) {
		return nil, true
	}
	// split file to header and footer and extend with license
	header, footer := splitSource(source, config.Golic.Rules[rule].Under)
	if header != "" {
		header = header + "\n"
	}
	source = fmt.Sprintf("%s%s%s", header, license, footer)

	if !o.Dry {
		data := []byte(source)
		err = os.WriteFile(path, data, os.ModeExclusive)
	}
	return
}

// removeFile Remove text from file overall function
func removeFile(path string, o internal.Options, config *Config) (err error, skip bool) {
	source, err := read(path)
	if err != nil {
		return err, false
	}
	rule := getRule(config, path)
	license, err := getCommentedLicense(config, o, rule)
	if err != nil {
		return err, false
	}
	if strings.Contains(source, license) {
		return RemoveFromFile(path, o, source, license, err), false
	}
	return nil, true
}

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
		return injectFile(path, o, config)
	case internal.LicenseRemove:
		return removeFile(path, o, config)
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

// read File
func read(f string) (s string, err error) {
	content, err := os.ReadFile(f)
	if err != nil {
		return
	}
	// Convert []byte to string and print to screen
	return string(content), nil
}

// RemoveFromFile
func RemoveFromFile(path string, o internal.Options, source string, license string, err error) error {
	if !o.Dry {
		source = strings.Replace(source, license, "", 1)
		err = os.WriteFile(path, []byte(source), os.ModeExclusive)
	}
	return err
}

// matchRule Match Rule
func matchRule(config *Config, fileName string) (rule string, ok bool) {
	if _, ok = config.Golic.Rules[fileName]; ok {
		return fileName, ok
	}
	for k := range config.Golic.Rules {
		if pkg.IsMatch(fileName, k) {
			return k, true
		}
	}
	return "", false
}

// getCommentedLicense Get Commented License File
func getCommentedLicense(config *Config, o internal.Options, file string) (string, error) {
	var ok bool
	var template string
	var rule string
	if template, ok = config.Golic.Licenses[o.Template]; !ok {
		return "", fmt.Errorf("no license found for %s, check configuration (.golic.yaml)", o.Template)
	}

	//if _, ok =  config.Golic.Rules[rule]; !ok {
	if rule, ok = matchRule(config, file); !ok {
		return "", fmt.Errorf("no rule found for %s, check configuration (.golic.yaml)", rule)
	}
	template = strings.ReplaceAll(template, "{{copyright}}", o.Copyright)
	if config.IsWrapped(rule) {
		return fmt.Sprintf("%s\n%s%s\n",
				config.Golic.Rules[rule].Prefix,
				template,
				config.Golic.Rules[rule].Suffix),
			nil
	}
	// `\r\n` -> `\r\n #`, `\n` -> `\n #`
	content := strings.ReplaceAll(template, "\n", fmt.Sprintf("\n%s", config.Golic.Rules[rule].Prefix))
	content = strings.TrimSuffix(content, config.Golic.Rules[rule].Prefix)
	content = config.Golic.Rules[rule].Prefix + content
	// "# \n" -> "#\n" // "# \r\n" -> "#\r\n"; some environments automatically remove spaces in empty lines. This makes problems in license PR's
	cleanedPrefix := strings.TrimSuffix(config.Golic.Rules[rule].Prefix, " ")
	content = strings.ReplaceAll(content, fmt.Sprintf("%s \n", cleanedPrefix), fmt.Sprintf("%s\n", cleanedPrefix))
	content = strings.ReplaceAll(content, fmt.Sprintf("%s \r\n", cleanedPrefix), fmt.Sprintf("%s\r\n", cleanedPrefix))
	return content, nil
}

// splitSource Split Source
func splitSource(source string, rules []string) (header, footer string) {
	lines := strings.Split(source, "\n")
	if len(rules) == 0 {
		return "", source
	}
	for _, r := range rules {
		header, footer = findHeaderAndFooter(lines, r)
		if header != "" {
			return
		}
	}
	return
}

func findHeaderAndFooter(lines []string, match string) (header, footer string) {
	for i, l := range lines {
		if pkg.IsMatch(l, match) {
			header = strings.Join(lines[0:i+1], "\n")
			footer = strings.Join(lines[i+1:], "\n")
			return
		}
	}
	return "", strings.Join(lines, "\n")
}

// getRule Get Rule for Match
func getRule(config *Config, path string) (rule string) {
	fileName := filepath.Base(path)
	for k := range config.Golic.Rules {
		matched, _ := filepath.Match(k, fileName)
		if matched {
			return k
		}
	}
	rule = filepath.Ext(path)
	if rule == "" {
		rule = filepath.Base(path)
	}
	return
}
