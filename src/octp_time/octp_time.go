package octp_time

import (
	"fmt"
	"time"
)

const (
	FORMAT_YEAR                 = 2006
	FORMAT_MONTH                = 1
	FORMAT_DAY                  = 2
	FORMAT_HOUR                 = 15
	FORMAT_MIN                  = 4
	FORMAT_SECOND               = 5
	FORMAT_DEFAULT_LONG_STRING  = "2006-01-02 15:04:05"
	FORMAT_DEFAULT_SHORT_STRING = "20060102150405"
)

func TimeCmper(ts1, ts2 interface{}) (rc int8, err error) {
	var k1, k2 string
	var ok bool
	if k1, ok = ts1.(string); !ok {
		err = fmt.Errorf("the argument is not string type")
		return
	}
	if k2, ok = ts2.(string); !ok {
		err = fmt.Errorf("the argument is not string type")
		return
	}
	var t1, t2 time.Time
	t1, err = time.Parse(FORMAT_DEFAULT_LONG_STRING, k1)
	if nil != err {
		return
	}
	t2, err = time.Parse(FORMAT_DEFAULT_LONG_STRING, k2)
	if nil != err {
		return
	}
	if t1.Equal(t2) {
		rc = 0
		return
	}
	if t1.Before(t2) {
		rc = -1
		return
	}
	if t1.After(t2) {
		rc = 1
		return
	}
	err = fmt.Errorf("can not compare time ts1:%q,ts2:%q.", ts1, ts2)
	rc = 0
	return
}

func GetNowShortString() (rc string) {
	now := time.Now()
	rc = now.Format(FORMAT_DEFAULT_SHORT_STRING)
	return
}
func GetNowLongString() (rc string) {
	now := time.Now()
	rc = now.Format(FORMAT_DEFAULT_LONG_STRING)
	return
}

func GetDateTimeShortString(t time.Time) (rc string) {
	rc = t.Format(FORMAT_DEFAULT_SHORT_STRING)
	return rc
}

func GetDateTimeLongString(t time.Time) (rc string) {
	rc = t.Format(FORMAT_DEFAULT_LONG_STRING)
	return rc
}
