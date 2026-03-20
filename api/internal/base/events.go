package base

import "time"

type Event struct {
	Timestamp time.Time `validate:"required"`
	AppId     string    `validate:"required"`
	KeyId     string    `validate:"required"`
	Cluster   string    `json:"cluster" validate:"max=256"` // required_if=EventType CSEV|required_if=EventType DCSEV,
	Tags      []string  `json:"tags" validate:"required"`
	TagString string
	Value     string `json:"value" validate:"max=256"` // required_if=EventType SEV|required_if=EventType CSEV,
}
