package types

import "github.com/navapbc/its-log/internal/constants"

type PermissionType int

func (ak *ApiKey) ConfigurePermissions() {
	switch ak.PermissionString {
	case "log":
		ak.Permission = constants.Log
	case "readonly":
		ak.Permission = constants.ReadOnly
	case "admin":
		ak.Permission = constants.Admin
	case "test":
		ak.Permission = constants.Test
	}

}

type ApiKey struct {
	AppId            string `json:"app_id" mapstructure:"app_id"`
	KeyId            string `json:"key_id" mapstructure:"key_id"`
	Key              string `json:"key" mapstructure:"key"`
	PermissionString string `json:"permission" mapstructure:"permission"`
	Permission       PermissionType
}

type ApiKeys []ApiKey
