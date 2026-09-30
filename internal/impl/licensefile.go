package impl

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
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	golog "github.com/Bugs5382/go-log"
	"github.com/Bugs5382/golic/internal/logging"
	"github.com/enescakir/emoji"
	"github.com/logrusorgru/aurora"
)

const defaultLicenseFile = "LICENSE"

// licenseTexts holds the full licence text for each built-in license, as
// downloaded from the source named in licenseFileSources. Never edit them.
//
//go:embed licenses/*.txt
var licenseTexts embed.FS

// licenseFileSource says where a LICENSE text comes from and which
// placeholder, if any, takes the copyright. Only texts that open with a
// copyright line for the licensee have one; the Apache and GNU texts carry
// their placeholders in a "how to apply" appendix, which stays as published.
type licenseFileSource struct {
	url         string
	placeholder string
}

var licenseFileSources = map[string]licenseFileSource{
	"apache2":   {url: "https://www.apache.org/licenses/LICENSE-2.0.txt"},
	"gpl2":      {url: "https://www.gnu.org/licenses/gpl-2.0.txt"},
	"gpl3":      {url: "https://www.gnu.org/licenses/gpl-3.0.txt"},
	"lgpl3":     {url: "https://www.gnu.org/licenses/lgpl-3.0.txt"},
	"agpl3":     {url: "https://www.gnu.org/licenses/agpl-3.0.txt"},
	"mpl2":      {url: "https://www.mozilla.org/media/MPL/2.0/index.txt"},
	"epl2":      {url: "https://www.eclipse.org/org/documents/epl-2.0/EPL-2.0.txt"},
	"unlicense": {url: "https://unlicense.org/UNLICENSE"},
	"mit":       {url: "https://github.com/spdx/license-list-data/blob/main/text/MIT.txt", placeholder: "<year> <copyright holders>"},
	"bsd2":      {url: "https://github.com/spdx/license-list-data/blob/main/text/BSD-2-Clause.txt", placeholder: "<year> <owner> "},
	"bsd3":      {url: "https://github.com/spdx/license-list-data/blob/main/text/BSD-3-Clause.txt", placeholder: "<year> <owner>. "},
	"isc":       {url: "https://github.com/spdx/license-list-data/blob/main/text/ISC.txt", placeholder: "<year> <owner> "},
}

// officialLicenseText returns the LICENSE text for a built-in license, with
// the copyright filled in where the text expects one.
func officialLicenseText(template, copyright string) (string, error) {
	src, ok := licenseFileSources[template]
	if !ok {
		keys := make([]string, 0, len(licenseFileSources))
		for k := range licenseFileSources {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return "", fmt.Errorf("no official LICENSE text for template %s; the LICENSE file option works with %s", template, strings.Join(keys, ", "))
	}
	raw, err := licenseTexts.ReadFile("licenses/" + template + ".txt")
	if err != nil {
		return "", err
	}
	text := string(raw)
	if src.placeholder != "" {
		text = strings.Replace(text, src.placeholder, copyright, 1)
	}
	return text, nil
}

// licenseFilePath is where the LICENSE goes for this run, or "" when the
// option is off. The flag wins over the config.
func (u *Process) licenseFilePath() string {
	if u.Opts.LicenseFile != "" {
		return u.Opts.LicenseFile
	}
	return u.cfg.Golic.LicenseFile.Path
}

// syncLicenseFile writes the official licence text to the LICENSE file when
// the option is on. A missing or different file counts as a change for -x.
func (u *Process) syncLicenseFile() error {
	path := u.licenseFilePath()
	if path == "" {
		logging.L().Trace("LICENSE file option off")
		return nil
	}
	logging.L().Debug("checking the LICENSE file", golog.F("path", path), golog.F("template", u.Opts.Template))

	want, err := officialLicenseText(u.Opts.Template, u.Opts.Copyright)
	if err != nil {
		return err
	}

	// #nosec G304 -- the path is the one the user named for the LICENSE file
	current, err := os.ReadFile(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		logging.L().Debug("LICENSE file missing", golog.F("path", path))
	case err != nil:
		return fmt.Errorf("failed to read %s: %w", path, err)
	case string(current) == want:
		logging.L().Info(fmt.Sprintf("%s  %s %s", emoji.Minus, aurora.Magenta(path), aurora.Bold(aurora.BrightMagenta("-> skip"))))
		return nil
	default:
		logging.L().Debug("LICENSE file differs from the official text", golog.F("path", path))
	}

	u.modified++
	prefix := ""
	if u.Opts.Dry {
		prefix = aurora.Bold(aurora.Yellow(fmt.Sprintf("%s DRY RUN: ", emoji.TestTube))).String()
	}
	logging.L().Info(fmt.Sprintf("%s %s  %s", prefix, emoji.Minus, aurora.BrightYellow(path)))
	if u.Opts.Dry {
		return nil
	}
	if err := os.WriteFile(path, []byte(want), 0o644); err != nil { // #nosec G306 -- a LICENSE file is meant to be world-readable
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	logging.L().Debug("LICENSE file written", golog.F("path", path), golog.F("bytes", len(want)))
	return nil
}
