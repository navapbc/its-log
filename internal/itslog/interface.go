package itslog

type ItsLog interface {
	Init() error
	Event(source string, event string, value any, value_type string) (int64, error)
}

type Event struct {
	Source string `json:"source" binding:"required"`
	Event  string `json:"event" binding:"required"`
	Value  string `json:"value" binding:"required"`
	Type   string `json:"type" binding:"required"`
}
