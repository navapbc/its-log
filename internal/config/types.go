package config

var KEY_KIND_LOGGING = "log"
var KEY_KIND_ADMIN = "admin"
var KEY_KIND_READONLY = "readonly"

type ApiKey struct {
	KeyID string `mapstructure:"id"`
	Key   string `mapstructure:"key"`
	Kind  string `mapstructure:"kind"`
}

// type ApiKeys struct {
// 	Keys []ApiKey `mapstructure:"api_keys"`
// }

type ApiKeys []ApiKey
