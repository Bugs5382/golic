package impl

type Config struct {
	Golic struct {
		Licenses map[string]string `yaml:"licenses"`
		Rules    map[string]struct {
			Prefix string   `yaml:"prefix"`
			Suffix string   `yaml:"suffix"`
			Under  []string `yaml:"under"`
		} `yaml:"rules"`
	} `yaml:"golic"`
}

func InitConfig() Config {
	cfg := Config{}
	return cfg
}

func (c *Config) IsWrapped(key string) bool {
	return c.Golic.Rules[key].Suffix != ""
}
