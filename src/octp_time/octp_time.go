package octp_time

import "time"

const (
	FORMAT_YEAR           = 2006
	FORMAT_MONTH          = 1
	FORMAT_DAY            = 2
	FORMAT_HOUR           = 15
	FORMAT_MIN            = 4
	FORMAT_SECOND         = 5
	FORMAT_DEFAULT_STRING = "2006-01-02 15:04:05"
)

func TimeCmper(ts1, ts2 string) uint8 {
	t1 := time.parser(FORMAT_DEFAULT_STRING, ts1)
	t2 := time.parser(FORMAT_DEFAULT_STRING, ts2)
	if t1.Equal(t2) {
		return 0
	}
	if t1.Before(t2) {
		return -1
	}
	if t1.After(t2) {
		return 1
	}
}
