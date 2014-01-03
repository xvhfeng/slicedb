package logger

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	path := "/opt/slicedb/log"
	prefix := "slicedb"
	level := 1
	var maxsize uint64 = 10 * 1024 * 1024
	/* var fileMode FileMode = os.ModePerm */
	log, err := New(path, level, prefix, maxsize)
	if nil != err {
		fmt.Println(err)
		return
	}
	/* for i := 0; i < 10*1024*1024; i++ { */
	err = log.Brush(Error, "today is fucking day,%d,%d", 900, 900)
	err = log.Brush(Error, "today is fucking day")
	err = log.Brush(Error, "today is fucking day")
	/* } */
}
