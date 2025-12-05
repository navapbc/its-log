package itslog

type ItsLog interface {
	Init() error
	Event(event string, value any, value_type string) (int64, error)
}

type Event struct {
	Application string `json:"app" binding:"required"`
	Event       string `json:"event" binding:"required"`
	Value       string `json:"value" binding:"required"`
	ValueType   string `json:"value_type" binding:"required"`
}
