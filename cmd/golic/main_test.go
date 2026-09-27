package main

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
	"os/exec"
	"path/filepath"
	"testing"
)

// buildBinary builds cmd/golic into a temp dir, optionally with -ldflags.
func buildBinary(t *testing.T, ldflags string) string {
	t.Helper()

	bin := filepath.Join(t.TempDir(), "golic")
	args := []string{"build", "-o", bin}
	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}
	args = append(args, ".")

	build := exec.Command("go", args...) // #nosec G204 -- fixed arguments, test only
	build.Env = append(build.Environ(), "CGO_ENABLED=0")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}
	return bin
}

// TestVersionLdflags builds the binary the way `task build` and GoReleaser do
// and checks that `golic version` prints the stamped values.
func TestVersionLdflags(t *testing.T) {
	if testing.Short() {
		t.Skip("builds the binary")
	}

	bin := buildBinary(t, "-X 'github.com/Bugs5382/golic/internal/build.Version=v9.8.7' "+
		"-X 'github.com/Bugs5382/golic/internal/build.Gitsha=abc1234'")

	out, err := exec.Command(bin, "version").Output() // #nosec G204 -- binary built above
	if err != nil {
		t.Fatalf("golic version failed: %v", err)
	}
	if got, want := string(out), "v9.8.7\ncommit abc1234\n"; got != want {
		t.Fatalf("golic version = %q, want %q", got, want)
	}
}
