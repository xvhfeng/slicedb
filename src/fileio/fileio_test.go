package fileio

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestBufferFile(t *testing.T) {
	f, err := os.OpenFile("/opt/slicedb/log.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if nil != err {
		fmt.Println("open bufio file is fail")
		return
	}

	b := NewBufWriter(f, 1024, time.Duration(10000*time.Second), nil)
	b.WriteString("1:i like\n")
	b.WriteString("2:i是中国人\n")
	b.WriteString("3:i like\n")
	/* time.Sleep(12 * time.Second) */
	b.WriteString("1:i like\n")
	b.WriteString("2:i是中国人\n")
	b.WriteString("3:i like\n")
	/* time.Sleep(11 * time.Second) */
	for i := 0; i < 1000000; i++ {
		b.WriteString("1:i like\n")
		b.WriteString("2:i是中国人\n")
		b.WriteString("3:i like\n")
		/* b.Write([]byte{byte(i)}) */
	}
	b.Flush()
	/* b.Flush() */
}
