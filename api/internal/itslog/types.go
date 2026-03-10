package itslog

import (
	"time"

	_ "github.com/creasty/defaults"
)

type Event struct {
	Timestamp time.Time `validate:"required"`
	AppId     string    `validate:"required"`
	KeyId     string    `validate:"required"`
	Cluster   string    `json:"cluster" validate:"max=256"` // required_if=EventType CSEV|required_if=EventType DCSEV,
	Tags      []string  `json:"tags" validate:"required"`
	TagString string
	Value     string `json:"value" validate:"max=256"` // required_if=EventType SEV|required_if=EventType CSEV,
}

type PermissionType int

const (
	Log PermissionType = iota + 1
	ReadOnly
	Admin
	Test
)

func (ak *ApiKey) ConfigurePermissions() {
	switch ak.PermissionString {
	case "log":
		ak.Permission = Log
	case "readonly":
		ak.Permission = ReadOnly
	case "admin":
		ak.Permission = Admin
	case "test":
		ak.Permission = Test
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

const OK = "ok"
const ERROR = "error"

type Success struct {
	Status string `json:"status"`
}
type Error struct {
	Status    string         `json:"status"`
	Data      map[string]any `json:"data"`
	ErrorType string         `json:"error_type"`
	Error     string         `json:"error"`
}
