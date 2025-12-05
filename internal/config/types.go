package config

var LOGGING_KEY_KIND = "log"
var ADMIN_KEY_KIND = "admin"

type ApiKey struct {
	Application string `mapstructure:"application"`
	Key         string `mapstructure:"key"`
	Kind        string `mapstructure:"kind"`
}

type ApiConfig struct {
	Keys []ApiKey `mapstructure:"api_keys"`
}
