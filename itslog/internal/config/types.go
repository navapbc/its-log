package config

var KEY_KIND_LOGGING = "log"
var KEY_KIND_ADMIN = "admin"
var KEY_KIND_READONLY = "readonly"
var KEY_KIND_TEST = "test"

type ApiKey struct {
	KeyId string `json:"key_id" mapstructure:"key_id"`
	Key   string `json:"key" mapstructure:"key"`
	Kind  string `json:"kind" mapstructure:"kind"`
}

// type ApiKeys struct {
// 	Keys []ApiKey `mapstructure:"api_keys"`
// }

type ApiKeys []ApiKey
