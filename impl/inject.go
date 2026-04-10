package impl

import (
	"context"

	"github.com/AbsaOSS/golic/internal"
	"github.com/denormal/go-gitignore"
	"github.com/enescakir/emoji"
	log "github.com/sirupsen/logrus"
)

func NewInject(ctx context.Context, options internal.Options) *Process {
	return &Process{
		Ctx:  ctx,
		Opts: options,

		modified: false,
	}
}

func (u *Process) Run() (err error) {
	// debug commands
	log.Debugf("%s reading config path: %s", emoji.OpenBook, u.Opts.ConfigPath)
	log.Debugf("%s reading lic ignore path: %s", emoji.OpenBook, u.Opts.LicIgnore)
	log.Debugf("%s reading template: %s", emoji.OpenBook, u.Opts.Template)
	log.Debugf("%s reading search path: %s", emoji.OpenBook, u.Opts.SearchPath)

	//if u.cfg, err = u.readCommonConfig(); err != nil {
	//	return
	//}
	// TODO Add Debug to output Config File Here... maybe

	u.ignore, err = gitignore.NewFromFile(u.Opts.LicIgnore)
	if err != nil {
		return err
	}
	// TODO Add Debug to output Ignore File Here... maybe

	return
}

func (u *Process) String() string {
	//TODO implement me
	panic("implement String()")
}

func (u *Process) ExitCode() int {
	if u.Opts.ModifiedExitStatus && u.modified {
		return 1
	}
	return 0
}

//
//func (u *InjectProcess) String() string {
//	//switch u.opts.Type {
//	//case pkg.LicenseInject:
//	//	return aurora.BrightCyan("inject").String()
//	//case pkg.LicenseRemove:
//	//	return aurora.BrightCyan("remove").String()
//	//}
//	//return aurora.BrightRed("ERROR, unrecognised command").String()
//}
//
//func (u *InjectProcess) ExitCode() int {
//	//if u.opts.ModifiedExitStatus && u.modified != 0 {
//	//	return 1
//	//}
//	//return 0
//}
//
//func read(f string) (s string, err error) {
//	//content, err := os.ReadFile(f)
//	//if err != nil {
//	//	return
//	//}
//	//// Convert []byte to string and print to screen
//	//return string(content), nil
//}
//
//func (u *InjectProcess) traverse() {
//	//skipped := 0
//	//visited := 0
//	//p := func(path string, i gitignore.GitIgnore, o pkg.Options, config *Config) (err error) {
//	//	if !i.Ignore(path) {
//	//		var skip bool
//	//		symbol := ""
//	//		cp := aurora.BrightYellow(path)
//	//		visited++
//	//		if err, skip = update(path, o, config); skip {
//	//			symbol = "-> skip"
//	//			cp = aurora.Magenta(path)
//	//			skipped++
//	//		}
//	//		_, _ = emoji.Printf(" %s  %s %s  \n", emoji.Minus, cp, aurora.BrightMagenta(symbol))
//	//	}
//	//	return
//	//}
//	//
//	//err := filepath.Walk("./",
//	//	func(path string, info os.FileInfo, err error) error {
//	//		if err != nil {
//	//			return err
//	//		}
//	//		if !info.IsDir() {
//	//			return p(path, u.ignore, u.opts, u.cfg)
//	//		}
//	//		return nil
//	//	})
//	//if err != nil {
//	//	logger.Err(err).Msg("")
//	//}
//	//u.modified = visited - skipped
//	//summary(skipped, visited)
//}
//
//func summary(skipped, visited int) {
//	//if skipped == visited {
//	//	fmt.Printf("\n %s %v/%v %s\n\n", emoji.Ice, aurora.BrightCyan(visited-skipped), aurora.BrightWhite(visited), aurora.BrightCyan("changed"))
//	//	return
//	//}
//	//fmt.Printf("\n %s %v/%v %s\n\n", emoji.Fire, aurora.BrightYellow(visited-skipped), aurora.BrightWhite(visited), aurora.BrightYellow("changed"))
//}
//
//func update(path string, o pkg.Options, config *Config) (err error, skip bool) {
//	//switch o.Type {
//	//case pkg.LicenseInject:
//	//	return inject(path, o, config)
//	//case pkg.LicenseRemove:
//	//	return remove(path, o, config)
//	//}
//	//return fmt.Errorf("invalid license type"), true
//}
//
//func inject(path string, o pkg.Options, config *Config) (err error, skip bool) {
//	//source, err := read(path)
//	//if err != nil {
//	//	return err, false
//	//}
//	//rule := getRule(config, path)
//	//license, err := getCommentedLicense(config, o, rule)
//	//if err != nil {
//	//	return err, false
//	//}
//	//// license is injected, continue
//	//if strings.Contains(source, license) {
//	//	return nil, true
//	//}
//	//// split file to header and footer and extend with license
//	//header, footer := splitSource(source, config.Golic.Rules[rule].Under)
//	//if header != "" {
//	//	header = header + "\n"
//	//}
//	//source = fmt.Sprintf("%s%s%s", header, license, footer)
//	//
//	//if !o.Dry {
//	//	data := []byte(source)
//	//	err = os.WriteFile(path, data, os.ModeExclusive)
//	//}
//	return
//}
//
//func remove(path string, o pkg.Options, config *Config) (err error, skip bool) {
//	//source, err := read(path)
//	//if err != nil {
//	//	return err, false
//	//}
//	//rule := getRule(config, path)
//	//license, err := getCommentedLicense(config, o, rule)
//	//if err != nil {
//	//	return err, false
//	//}
//	//if strings.Contains(source, license) {
//	//	return RemoveFromFile(path, o, source, license, err), false
//	//}
//	return nil, true
//}
//
//func RemoveFromFile(path string, o pkg.Options, source string, license string, err error) error {
///*	if !o.Dry {
//		source = strings.Replace(source, license, "", 1)
//		err = os.WriteFile(path, []byte(source), os.ModeExclusive)
//	}*/
//	return err
//}
//
//func matchRule(config *Config, fileName string) (rule string, ok bool) {
///*	if _, ok = config.Golic.Rules[fileName]; ok {
//		return fileName, ok
//	}
//	for k := range config.Golic.Rules {
//		if update2.isMatch(fileName, k) {
//			return k, true
//		}
//	}*/
//	return "", false
//}
//
//func getCommentedLicense(config *Config, o pkg.Options, file string) (string, error) {
//	//var ok bool
//	//var template string
//	//var rule string
//	//if template, ok = config.Golic.Licenses[o.Template]; !ok {
//	//	return "", fmt.Errorf("no license found for %s, check configuration (.golic.yaml)", o.Template)
//	//}
//	//
//	//log.Debugf("template: %file", file)
//	//
//	////if _, ok =  config.Golic.Rules[rule]; !ok {
//	//if rule, ok = matchRule(config, file); !ok {
//	//	return "", fmt.Errorf("no rule found for %s, check configuration (.golic.yaml)", rule)
//	//}
//	//template = strings.ReplaceAll(template, "{{copyright}}", o.Copyright)
//	//if config.IsWrapped(rule) {
//	//	return fmt.Sprintf("%s\n%s%s\n",
//	//			config.Golic.Rules[rule].Prefix,
//	//			template,
//	//			config.Golic.Rules[rule].Suffix),
//	//		nil
//	//}
//	//// `\r\n` -> `\r\n #`, `\n` -> `\n #`
//	//content := strings.ReplaceAll(template, "\n", fmt.Sprintf("\n%s", config.Golic.Rules[rule].Prefix))
//	//content = strings.TrimSuffix(content, config.Golic.Rules[rule].Prefix)
//	//content = config.Golic.Rules[rule].Prefix + content
//	//// "# \n" -> "#\n" // "# \r\n" -> "#\r\n"; some environments automatically remove spaces in empty lines. This makes problems in license PR's
//	//cleanedPrefix := strings.TrimSuffix(config.Golic.Rules[rule].Prefix, " ")
//	//content = strings.ReplaceAll(content, fmt.Sprintf("%s \n", cleanedPrefix), fmt.Sprintf("%s\n", cleanedPrefix))
//	//content = strings.ReplaceAll(content, fmt.Sprintf("%s \r\n", cleanedPrefix), fmt.Sprintf("%s\r\n", cleanedPrefix))
//	//return content, nil
//}
//
//func splitSource(source string, rules []string) (header, footer string) {
//	lines := strings.Split(source, "\n")
//	if len(rules) == 0 {
//		return "", source
//	}
//	for _, r := range rules {
//		header, footer = findHeaderAndFooter(lines, r)
//		if header != "" {
//			return
//		}
//	}
//	return
//}
//
//func findHeaderAndFooter(lines []string, match string) (header, footer string) {
//	for i, l := range lines {
//		if update2.isMatch(l, match) {
//			header = strings.Join(lines[0:i+1], "\n")
//			footer = strings.Join(lines[i+1:], "\n")
//			return
//		}
//	}
//	return "", strings.Join(lines, "\n")
//}
//
//func getRule(config *Config, path string) (rule string) {
//	fileName := filepath.Base(path)
//	for k := range config.Golic.Rules {
//		matched, _ := filepath.Match(k, fileName)
//		if matched {
//			return k
//		}
//	}
//	rule = filepath.Ext(path)
//	if rule == "" {
//		rule = filepath.Base(path)
//	}
//	return
//}
//
//func (u *InjectProcess) readLocalConfig() (*Config, error) {
//	var c = &Config{}
//	var rc = *u.cfg
//	yamlFile, err := os.ReadFile(u.opts.ConfigPath)
//	if err != nil {
//		return nil, nil
//	}
//	err = yaml.Unmarshal(yamlFile, c)
//	if err != nil {
//		return nil, nil
//	}
//	for k, v := range c.Golic.Licenses {
//		rc.Golic.Licenses[k] = v
//	}
//	for k, v := range c.Golic.Rules {
//		rc.Golic.Rules[k] = v
//	}
//	return &rc, nil
//}
//
//func (u *InjectProcess) readCommonConfig() (c *Config, err error) {
//	c = &Config{}
//	err = yaml.Unmarshal([]byte(u.opts.MasterConfig), c)
//	return
//}
