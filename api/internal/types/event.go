package types

import "time"

type Event struct {
	Timestamp time.Time
	AppId     string
	KeyId     string
	Cluster   string   `json:"cluster"`
	Tags      []string `json:"tags"`
	TagString string
	Value     string `json:"value"`
	Date      string `json:"date"`
}
