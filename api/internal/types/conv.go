package types

import (
	"fmt"
	"time"
)

func YMDToUnixEpoch(ymd string) (int64, error) {
	timetime, err := time.Parse("2006-01-02", ymd)
	if err != nil {
		return -1, fmt.Errorf("YMDToUnixEpoch: %s", err.Error())
	}
	return timetime.Unix(), nil
}
