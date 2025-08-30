package main

type AppConfig struct {
	Domain string `mapstructure:"domain"`
	Type   string `mapstructure:"type"`
	Port   int    `mapstructure:"port"`
	Path   string `mapstructure:"path"`
}

type Config struct {
	App map[string]AppConfig `mapstructure:"app"`
}
