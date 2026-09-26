package build

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
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

// Version and Gitsha are stamped at link time by `task build` and GoReleaser:
//
//	-X 'github.com/Bugs5382/golic/internal/build.Version=v1.0.0'
//	-X 'github.com/Bugs5382/golic/internal/build.Gitsha=abc1234'
//
// They stay empty for `go install` and plain `go build`, where the version
// comes from the module build info instead.
var (
	Version = ""
	Gitsha  = ""
)

// Devel is reported when neither the linker nor the build info carries a
// release version, for example a `go run` or a build outside a module.
const Devel = "(devel)"

// shortCommitLen matches `git rev-parse --short` and the GoReleaser
// ShortCommit template.
const shortCommitLen = 7

// Info returns the version and commit of the running binary. The linker
// values win; otherwise the module version and vcs revision recorded by the
// Go toolchain are used. The commit is empty when it is not known.
func Info() (version, commit string) {
	bi, ok := debug.ReadBuildInfo()
	return resolve(Version, Gitsha, bi, ok)
}

// String renders Info as the `golic version` output: the version on the first
// line and, when known, the commit on a second line.
func String() string {
	return format(Info())
}

func format(version, commit string) string {
	if commit == "" {
		return version + "\n"
	}
	return fmt.Sprintf("%s\ncommit %s\n", version, commit)
}

func resolve(ldVersion, ldCommit string, bi *debug.BuildInfo, ok bool) (version, commit string) {
	version, commit = ldVersion, ldCommit

	if version != "" {
		log.Debug().Str("version", version).Msg("version taken from the linker flags")
	}

	if ok && bi != nil {
		if version == "" {
			switch bi.Main.Version {
			case "", Devel:
				log.Debug().Str("module_version", bi.Main.Version).Msg("build info has no release version")
			default:
				version = bi.Main.Version
				log.Debug().Str("version", version).Msg("version taken from the module build info")
			}
		}
		if commit == "" {
			commit = vcsRevision(bi)
		}
	} else {
		log.Debug().Msg("no build info available")
	}

	if version == "" {
		version = Devel
	}
	log.Trace().Str("version", version).Str("commit", commit).Msg("version resolved")
	return version, commit
}

func vcsRevision(bi *debug.BuildInfo) string {
	for _, s := range bi.Settings {
		if s.Key != "vcs.revision" {
			continue
		}
		rev := s.Value
		if len(rev) > shortCommitLen {
			rev = rev[:shortCommitLen]
		}
		log.Trace().Str("commit", rev).Msg("commit taken from the vcs build setting")
		return rev
	}
	return ""
}
