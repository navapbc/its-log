package itslog

type ItsLog interface {
	Init() error
	Event(event string, value any, value_type LogItType) (int64, error)
}

type Event struct {
	Event     string `json:"event" binding:"required"`
	Value     string `json:"value" binding:"required"`
	ValueType int64  `json:"value_type" binding:"required"`
}
