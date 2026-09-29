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
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bugs5382/golic/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// builtinIgnores is the default ignore list from #46. The embedded config must
// ship exactly these, in this order.
var builtinIgnores = []string{
	".git/**",
	"**/node_modules/**",
	"**/vendor/**",
	"**/dist/**",
	"**/coverage/**",
	"**/*.pb.go",
	"**/*_gen.go",
	"**/*.gen.go",
	"**/zz_generated*.go",
	"**/*.connect.go",
	"**/*.d.ts",
	"**/*.lock",
	"**/package-lock.json",
	"**/go.sum",
}

func TestEmbeddedConfigShipsDefaultIgnores(t *testing.T) {
	assert.Equal(t, builtinIgnores, embeddedConfig(t).Golic.Ignore)
}

// TestDefaultIgnorePatterns checks every built-in pattern against paths it
// must take out of scope, at the root and nested, and near misses it must not.
func TestDefaultIgnorePatterns(t *testing.T) {
	cases := []struct {
		pattern string
		ignored []string
		kept    []string
	}{
		{".git/**", []string{".git/config", ".git/hooks/pre-commit.sh"}, []string{"git/main.go", "a.git.go"}},
		{"**/node_modules/**", []string{"node_modules/x/index.js", "web/node_modules/y/a.ts"}, []string{"node_modules.go", "src/modules/a.ts"}},
		{"**/vendor/**", []string{"vendor/github.com/x/y.go", "svc/vendor/z.go"}, []string{"vendored.go", "src/vendors/a.go"}},
		{"**/dist/**", []string{"dist/index.js", "packages/a/dist/b.mjs"}, []string{"distance.go", "src/distro/a.go"}},
		{"**/coverage/**", []string{"coverage/lcov-report/a.js", "pkg/coverage/b.js"}, []string{"coverage.go", "src/cover/a.go"}},
		{"**/*.pb.go", []string{"x.pb.go", "proto/livenesspb/liveness.pb.go"}, []string{"pb.go", "proto/liveness.go"}},
		{"**/*_gen.go", []string{"enum_gen.go", "internal/x/types_gen.go"}, []string{"gen.go", "generate.go"}},
		{"**/*.gen.go", []string{"api.gen.go", "internal/api/server.gen.go"}, []string{"gen.go", "api.go"}},
		{"**/zz_generated*.go", []string{"zz_generated.deepcopy.go", "api/v1/zz_generated_conversion.go"}, []string{"generated.go", "api/v1/types.go"}},
		{"**/*.connect.go", []string{"foo.connect.go", "gen/foov1connect/foo.connect.go"}, []string{"connect.go", "server/connect_test.go"}},
		{"**/*.d.ts", []string{"index.d.ts", "types/api.d.ts"}, []string{"index.ts", "src/d.ts.go"}},
		{"**/*.lock", []string{"yarn.lock", "rust/Cargo.lock"}, []string{"lock.go", "lockfile.yaml"}},
		{"**/package-lock.json", []string{"package-lock.json", "web/package-lock.json"}, []string{"package.json"}},
		{"**/go.sum", []string{"go.sum", "tools/go.sum"}, []string{"go.mod", "go.sum.txt"}},
	}
	require.Len(t, cases, len(builtinIgnores), "one case per built-in pattern")
	for _, c := range cases {
		t.Run(c.pattern, func(t *testing.T) {
			for _, p := range c.ignored {
				got, ok := matchIgnore([]string{c.pattern}, p)
				assert.True(t, ok, "%s must ignore %s", c.pattern, p)
				assert.Equal(t, c.pattern, got)
				_, ok = matchIgnore(builtinIgnores, p)
				assert.True(t, ok, "the built-in list must ignore %s", p)
			}
			for _, p := range c.kept {
				_, ok := matchIgnore(builtinIgnores, p)
				assert.False(t, ok, "the built-in list must not ignore %s", p)
			}
		})
	}
}

// TestCatchAllNegations pins which .licignore negations are too broad to lift
// a built-in ignore (whitelist lines) and which name something specific.
func TestCatchAllNegations(t *testing.T) {
	for _, p := range []string{"*", "*/", "**", "/*", "*.go", "**/*.go", "*.ts", "*.lock", "/*.yaml"} {
		assert.True(t, isCatchAll(p, builtinIgnores), "%q is a catch-all", p)
	}
	for _, p := range []string{"*.pb.go", "**/*.pb.go", "**/*.lock", "vendor/**", "yarn.lock", "proto/x.pb.go", "*.d.ts", "Makefile"} {
		assert.False(t, isCatchAll(p, builtinIgnores), "%q is specific", p)
	}
}

// runTree writes files (including any .licignore or .golic.yaml among them)
// into a fresh directory, runs an inject pass from its root the way the CLI
// does, and returns every file's content afterwards.
func runTree(t *testing.T, files map[string]string) map[string]string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o750))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o600))
	}
	t.Chdir(dir)

	opts := internal.Options{
		Type:      internal.LicenseInject,
		Template:  "apache2",
		Copyright: testCopyright,
		LicIgnore: ".licignore",
	}
	if _, ok := files[".golic.yaml"]; ok {
		opts.ConfigPath = ".golic.yaml"
	}
	u := ProcessFile(context.Background(), opts)
	require.NoError(t, u.Run())

	out := make(map[string]string, len(files))
	for name := range files {
		b, err := os.ReadFile(filepath.FromSlash(name))
		require.NoError(t, err)
		out[name] = string(b)
	}
	return out
}

const plainGo = "package x\n\nfunc F() {}\n"

func stamped(t *testing.T, path, content string) bool {
	t.Helper()
	return strings.Contains(content, header(t, path, "apache2"))
}

func TestDefaultIgnoresKeepFilesOutOfScope(t *testing.T) {
	files := map[string]string{
		".licignore":                       "",
		"main.go":                          plainGo,
		"vendor/github.com/a/b.go":         plainGo,
		"node_modules/pkg/index.js":        "module.exports = 1\n",
		"web/dist/app.js":                  "console.log(1)\n",
		"coverage/lcov-report/prettify.js": "var a = 1\n",
		"proto/livenesspb/liveness.pb.go":  plainGo,
		"internal/enum_gen.go":             plainGo,
		"internal/api/server.gen.go":       plainGo,
		"api/v1/zz_generated.deepcopy.go":  plainGo,
		"gen/foov1connect/foo.connect.go":  plainGo,
		"types/index.d.ts":                 "export {}\n",
		"src/index.ts":                     "export {}\n",
	}
	out := runTree(t, files)
	for name := range files {
		switch name {
		case ".licignore":
		case "main.go", "src/index.ts":
			assert.True(t, stamped(t, name, out[name]), "%s is in scope and gets a header", name)
		default:
			assert.Equal(t, files[name], out[name], "%s is ignored by default and must not change", name)
		}
	}
}

// TestWhitelistLicignoreDoesNotLiftDefaults is the common "ignore *, allow
// *.go and */" .licignore: its broad negations must not pull vendored or
// generated files back in.
func TestWhitelistLicignoreDoesNotLiftDefaults(t *testing.T) {
	files := map[string]string{
		".licignore":                      "*\n!*.go\n!*/\n",
		"main.go":                         plainGo,
		"vendor/github.com/a/b.go":        plainGo,
		"proto/livenesspb/liveness.pb.go": plainGo,
	}
	out := runTree(t, files)
	assert.True(t, stamped(t, "main.go", out["main.go"]))
	assert.Equal(t, plainGo, out["vendor/github.com/a/b.go"])
	assert.Equal(t, plainGo, out["proto/livenesspb/liveness.pb.go"])
}

// TestLicignoreReincludeOverridesDefault checks a specific "!" line in
// .licignore lifts a built-in ignore, and only for the paths it names.
func TestLicignoreReincludeOverridesDefault(t *testing.T) {
	files := map[string]string{
		".licignore":           "*\n!*.go\n!*/\n!**/*.pb.go\n!vendor/ours/**\n",
		"api/x.pb.go":          plainGo,
		"vendor/ours/a.go":     plainGo,
		"vendor/theirs/b.go":   plainGo,
		"internal/enum_gen.go": plainGo,
	}
	out := runTree(t, files)
	assert.True(t, stamped(t, "api/x.pb.go", out["api/x.pb.go"]), "!**/*.pb.go re-includes pb files")
	assert.True(t, stamped(t, "vendor/ours/a.go", out["vendor/ours/a.go"]), "!vendor/ours/** re-includes that folder")
	assert.Equal(t, plainGo, out["vendor/theirs/b.go"], "the rest of vendor stays ignored")
	assert.Equal(t, plainGo, out["internal/enum_gen.go"], "other defaults still apply")
}

// TestLicignoreStillIgnoresAfterReinclude checks .licignore keeps the last
// word: a later ignore line wins over an earlier specific re-include.
func TestLicignoreStillIgnoresAfterReinclude(t *testing.T) {
	files := map[string]string{
		".licignore":  "!**/*.pb.go\napi/**\n",
		"api/x.pb.go": plainGo,
	}
	out := runTree(t, files)
	assert.Equal(t, plainGo, out["api/x.pb.go"])
}

// TestGitDirIsNeverWalked makes .git unreadable: walking into it would fail
// the run. Even a .licignore that re-includes it cannot pull it back in.
func TestGitDirIsNeverWalked(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".git", "objects"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".git", "hook.sh"), []byte("echo hi\n"), 0o600))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "sub", ".git"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "sub", ".git", "x.yaml"), []byte("a: 1\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".licignore"), []byte("!.git/**\n!**/.git/**\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.go"), []byte(plainGo), 0o600))
	require.NoError(t, os.Chmod(filepath.Join(dir, ".git", "objects"), 0o000))
	t.Cleanup(func() { _ = os.Chmod(filepath.Join(dir, ".git", "objects"), 0o750) })
	t.Chdir(dir)

	u := ProcessFile(context.Background(), internal.Options{
		Type: internal.LicenseInject, Template: "apache2", Copyright: testCopyright, LicIgnore: ".licignore",
	})
	require.NoError(t, u.Run(), "the walk must not enter .git")

	b, err := os.ReadFile(filepath.Join(dir, ".git", "hook.sh"))
	require.NoError(t, err)
	assert.Equal(t, "echo hi\n", string(b))
	b, err = os.ReadFile(filepath.Join(dir, "sub", ".git", "x.yaml"))
	require.NoError(t, err)
	assert.Equal(t, "a: 1\n", string(b), "nested .git folders are skipped too")
	assert.Equal(t, 2, u.modified, "only main.go and .licignore are in scope")
}

// TestLocalConfigIgnoreAddsToDefaults checks ignore: entries in a local
// .golic.yaml are added to the built-in list, not swapped for it.
func TestLocalConfigIgnoreAddsToDefaults(t *testing.T) {
	files := map[string]string{
		".licignore":               "*\n!*.go\n!*/\n",
		".golic.yaml":              "golic:\n  ignore:\n    - \"third_party/**\"\n",
		"main.go":                  plainGo,
		"third_party/x/a.go":       plainGo,
		"vendor/github.com/a/b.go": plainGo,
	}
	out := runTree(t, files)
	assert.True(t, stamped(t, "main.go", out["main.go"]))
	assert.Equal(t, plainGo, out["third_party/x/a.go"])
	assert.Equal(t, plainGo, out["vendor/github.com/a/b.go"])
}

func TestGeneratedFilesAreSkippedByContent(t *testing.T) {
	gen := "// Code generated by sqlc. DO NOT EDIT.\n// versions:\n//   sqlc v1.27.0\n\npackage db\n"
	files := map[string]string{
		".licignore":   "",
		"db/query.go":  gen,
		"db/manual.go": plainGo,
	}
	out := runTree(t, files)
	assert.Equal(t, gen, out["db/query.go"], "a file with a generated marker is left alone")
	assert.True(t, stamped(t, "db/manual.go", out["db/manual.go"]))
}
