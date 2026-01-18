package itslog

import (
	"fmt"
	"slices"
	"time"

	_ "github.com/creasty/defaults"
)

type EventType int

// Starting a 1 because the validation library
// treats zero integer values as missing.
const (
	SE = iota + 1
	SEV
	CSE
	CSEV
	DSE
	DSEV
	DCSE
	DCSEV
)

func (e EventType) Validate() error {
	valid := []EventType{SE, SEV, CSE, CSEV, DSE, DSEV, DCSE, DCSEV}
	if slices.Contains(valid, e) {
		return nil
	}
	return fmt.Errorf("%v unknown EventType", e)
}

type Event struct {
	Timestamp time.Time `validate:"required"`
	AppId     string    `validate:"required"`
	KeyId     string    `validate:"required"`
	EventType EventType `validate:"required,validateFn"`
	Cluster   string    `validate:"max=256"` // required_if=EventType CSEV|required_if=EventType DCSEV,
	Source    string    `validate:"required,max=256"`
	Event     string    `validate:"required,max=256"`
	Value     string    `validate:"max=256"` // required_if=EventType SEV|required_if=EventType CSEV,
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
