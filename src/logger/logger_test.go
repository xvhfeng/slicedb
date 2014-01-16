package logger

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	/* path := "/opt/slicedb/log" */
	/* prefix := "slicedb" */
	/* level := 1 */
	/* var maxsize uint64 = 10 * 1024 * 1024 */
	/* var fileMode FileMode = os.ModePerm */
	/* log, err := New(path, level, prefix, maxsize) */
	/* if nil != err { */
	/* fmt.Println(err) */
	/* return */
	/* } */
	/* for i := 0; i < 10*1024*1024; i++ { */
	/* err = log.Brush(Error, "today is fucking day,%d,%d", 900, 900) */
	/* err = log.Brush(Error, "today is fucking day") */
	/* err = log.Brush(Error, "today is fucking day") */
	/* var f *os.File */
	f, err := os.OpenFile("/opt/slicedb/log.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if nil != err {
		fmt.Println("open bufio file is fail")
		return
	}
	var w *bufio.Writer
	w = bufio.NewWriter(f)
	_, err = w.WriteString("我是中国人")
	if nil != err {
		fmt.Println("write chines is fail.")
	}
	w.Flush()
	f.Close()

	/* } */
}
