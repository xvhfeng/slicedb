/*
* binlog format:op|file-idx|offset|len|
 */
package binlog

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	W = 0x01
	D = 0x02
)

type Writer struct {
	locker   sync.Mutex
	basepath string
	t        time.Time
	timeSpan int64
	idx      int
	maxsize  uint64
	size     uint64
	f        *os.File
	log      *logger.Log
}

func binlogFilename(t time.Time, basepath string, idx int) (filename string) {
	s := octp_time.GetDateTimeShortString(t)
	filename := strings.Join([]string{basepath, s,
		".binlog_", strconv.itoa(idx)}, "")
	return filename
}

func NewWriter(basepath string, time time.Time,
	timeSpan int64, maxsize uint64,
	log *logger.Log) (binlog *Writer, err error) {
	if 0 == len(basepath) {
		err = fmt.Errorf("the argument:basepath is error.")
		return
	}
	err = os.MkdirAll(basepath, 0777)
	if nil != err {
		if os.IsExist(err) {
			err = nil
		} else {
			return
		}
	}
	if !strings.HasSuffix(path, "/") {
		basepath = strings.Join([]string{basepath, string(os.PathSeparator)}, "")
	}
	binlog = new(Writer)
	binlog.basepath = basepath
	binlog.maxsize = maxsize
	binlog.t = time
	binlog.timespan = timeSpan
	binlog.log = log
	err = binlog.open()
	return
}

func (binlog *Writer) open() (err error) {
	t := time.Now()
	if t.Sub(binlog.t).Seconds() > binlog.timespan {
		binlog.t = t
	}
	p := binlogFilename(binlog.t, binlog.basepath, binlog.idx)
	binlog.f, err = os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	return
}

func (binlog *Writer) Close() (err error) {
	if nil != binlog.f {
		err = binlog.f.Close()
	}
	binlog.size = 0
	return
}

func (binlog *Writer) Brush(format string, args ...interface{}) (err error) {
	line := fmt.Sprintf(format, args...)
	binlog.locker.Lock()
	defer binlog.locker.Unlock()
	if binlog.maxsize <= binlog.size {
		binlog.Close()
		binlog.idx++
		binlog.open()
	}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, line)
	err = binary.Write(buf, binary.LittleEndian, '\n')
	if nil != err {
		if nil != binlog.log {
			binlog.log.Brush(logger.Error, "write binlog:%v is fail.",
				line)
		}
		return err
	}
	bs := buf.Bytes()
	binlog.f.Write(bs)
	if nil == err {
		binlog.size += uint64(len(bs))
	}
	return
}

type Reader struct {
	offset int64
	f      *os.File
	reader bufio.Reader
	log    *logger.Log
}

func NewReader(path string, t time.Time, idx int,
	offset int64, log *logger.Log) (Reader, error) {
	var r *Reader

	r = new(Reader)
	fn := binlogFilename(t, path, idx)
	if f, err := os.Open(fn); nil != err {
		return nil, err
	} else {
		if 0 != offset {
			if _, err = f.Seek(offset, os.SEEK_CUR); nil != err {
				return nil, err
			}
		}
		r.f = f
		r.reader = bufio.NewReader(f)
		r.offset = offset
		r.log = log
	}
	return r, nil
}

func (r *Reader) Read() (line []byte, err error) {
	if line, err = r.reader.ReadBytes('\n'); nil != err {
		return
	}
	r.offset += len(line)
	return
}

func (r *Reader) ExpLine(buff []byte) (line string, err error) {
	b := bytes.NewReader(buff)
	err = binary.Read(b, binary.LittleEndian, &line)
	return
}

func (r *Reader) Close() {
	r.f.Close()
}
