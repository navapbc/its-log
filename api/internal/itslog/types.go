// Package itslog implements core types and operations for the its-log logger
//
// Amongst the key types and functions provided in this package are:
//
// EventType, which provides an enumeration of types of events that can be logged to its-log
// Event, the struct for events that are logged by the API
// ApiKey, which structures information around API keys (including app and key IDs)
package itslog

import (
	"fmt"
	"slices"
	"time"
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
	Logging PermissionType = iota
	ReadOnly
	Admin
	Test
)

func (ak *ApiKey) ConfigurePermissions() {
	switch ak.PermissionString {
	case "logging":
		ak.Permission = Logging
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
