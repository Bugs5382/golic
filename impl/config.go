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

type Rule struct {
	Prefix string   `yaml:"prefix"`
	Suffix string   `yaml:"suffix"`
	Under  []string `yaml:"under"`
}

type GolicConfig struct {
	Licenses   map[string]string `yaml:"licenses"`
	MergeRules bool              `yaml:"mergeRules"`
	Rules      map[string]Rule   `yaml:"rules"`
}

type Config struct {
	Golic GolicConfig `yaml:"golic"`
}

func (c *Config) IsWrapped(key string) bool {
	return c.Golic.Rules[key].Suffix != ""
}
