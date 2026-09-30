package build

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
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildInfo(version string, settings ...debug.BuildSetting) *debug.BuildInfo {
	return &debug.BuildInfo{
		Main:     debug.Module{Path: "github.com/Bugs5382/golic", Version: version},
		Settings: settings,
	}
}

func TestResolve(t *testing.T) {
	tests := []struct {
		name        string
		ldVersion   string
		ldCommit    string
		info        *debug.BuildInfo
		ok          bool
		wantVersion string
		wantCommit  string
	}{
		{
			name:        "ldflags value wins over build info",
			ldVersion:   "v1.2.3",
			ldCommit:    "abc1234",
			info:        buildInfo("v0.0.9"),
			ok:          true,
			wantVersion: "v1.2.3",
			wantCommit:  "abc1234",
		},
		{
			name:        "ldflags version without a commit falls back to the vcs revision",
			ldVersion:   "v1.2.3",
			info:        buildInfo("(devel)", debug.BuildSetting{Key: "vcs.revision", Value: "0123456789abcdef0123456789abcdef01234567"}),
			ok:          true,
			wantVersion: "v1.2.3",
			wantCommit:  "0123456",
		},
		{
			name:        "go install module version from build info",
			info:        buildInfo("v1.0.0"),
			ok:          true,
			wantVersion: "v1.0.0",
			wantCommit:  "",
		},
		{
			name:        "local build with a pseudo-version and vcs revision",
			info:        buildInfo("v0.3.1-0.20260926210015-14eb84d5eb7f+dirty", debug.BuildSetting{Key: "vcs.revision", Value: "14eb84d5eb7fe99ef6e9c529447d7bef2a36304a"}),
			ok:          true,
			wantVersion: "v0.3.1-0.20260926210015-14eb84d5eb7f+dirty",
			wantCommit:  "14eb84d",
		},
		{
			name:        "devel build reports devel",
			info:        buildInfo("(devel)"),
			ok:          true,
			wantVersion: Devel,
			wantCommit:  "",
		},
		{
			name:        "empty module version reports devel",
			info:        buildInfo(""),
			ok:          true,
			wantVersion: Devel,
			wantCommit:  "",
		},
		{
			name:        "no build info reports devel",
			info:        nil,
			ok:          false,
			wantVersion: Devel,
			wantCommit:  "",
		},
		{
			name:        "short vcs revision is kept as is",
			info:        buildInfo("(devel)", debug.BuildSetting{Key: "vcs.revision", Value: "abc"}),
			ok:          true,
			wantVersion: Devel,
			wantCommit:  "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersion, gotCommit := resolve(tt.ldVersion, tt.ldCommit, tt.info, tt.ok)
			assert.Equal(t, tt.wantVersion, gotVersion)
			assert.Equal(t, tt.wantCommit, gotCommit)
		})
	}
}

func TestString(t *testing.T) {
	assert.Equal(t, "v1.0.0\n", format("v1.0.0", ""))
	assert.Equal(t, "v1.0.0\ncommit abc1234\n", format("v1.0.0", "abc1234"))
}

func TestInfoUsesTheLinkedVariables(t *testing.T) {
	oldVersion, oldCommit := Version, Gitsha
	t.Cleanup(func() { Version, Gitsha = oldVersion, oldCommit })

	Version, Gitsha = "v9.9.9", "fedcba9"

	version, commit := Info()
	assert.Equal(t, "v9.9.9", version)
	assert.Equal(t, "fedcba9", commit)
	assert.Equal(t, "v9.9.9\ncommit fedcba9\n", String())
}
