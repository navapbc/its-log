package types

import (
	"time"
)

type ILTime struct {
	theTime time.Time
}

func NewILTimeNow() *ILTime {
	return &ILTime{theTime: time.Now()}
}

func NewILTimeToday() *ILTime {
	t := time.Now()
	return &ILTime{theTime: truncateToDay(t)}
}

func ILTimeFromTime(t time.Time) *ILTime {
	return &ILTime{theTime: t}
}

func ILTimeFromYMD(ymd string) (*ILTime, error) {
	d, e := time.Parse("2006-01-02", ymd)
	if e != nil {
		return nil, e
	}
	ilt := ILTimeFromTime(d)
	return ilt, nil
}

func (ilt *ILTime) offsetDays(days int) {
	ilt.theTime = ilt.theTime.AddDate(0, 0, days)
}

func (ilt *ILTime) SubtractDays(days int) {
	offset := -1 * days
	// DEBUG LOG
	// log.Printf("SubtractDays: %d\n", offset)
	ilt.offsetDays(offset)
}

func (ilt *ILTime) AddDays(days int) {
	ilt.offsetDays(days)
}

func (ilt *ILTime) AsTime() time.Time {
	return ilt.theTime
}

func (ilt *ILTime) AsYYYYMMDD() string {
	ymd := ilt.theTime.Format("2006-01-02")
	return ymd
}

func (ilt *ILTime) AsEpoch() int64 {
	return ilt.theTime.Unix()
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
