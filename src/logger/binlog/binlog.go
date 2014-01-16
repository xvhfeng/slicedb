package binlog

import (
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

type Binlog struct {
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

func NewBinlog(basepath string, time time.Time,
	timeSpan int64, maxsize uint64,
	log *logger.Log) (binlog *Binlog, err error) {
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
	binlog = new(Binlog)
	binlog.basepath = basepath
	binlog.maxsize = maxsize
	binlog.t = time
	binlog.timespan = timeSpan
	binlog.log = log
	err = binlog.open()
	return
}

func (binlog *Binlog) open() (err error) {
	t := time.Now()
	if t.Sub(binlog.t).Seconds() > binlog.timespan {
		binlog.t = t
	}
	p := binlogFilename(binlog.t, binlog.basepath, binlog.idx)
	binlog.f, err = os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	return
}

func (binlog *Binlog) Close() (err error) {
	if nil != binlog.f {
		err = binlog.f.Close()
	}
	binlog.size = 0
	return
}

func (binlog *Binlog) Brush(format string, args ...interface{}) (err error) {
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

/* const BINLOG_READER_BUFFER_SIZE = 1024 */
/*
type BinlogReader struct {
    basepath string
    t        time.Time
    timeSpan int64
    idx      int
    offset int64
    maxsize	 int64
    f        *os.File
    log      *logger.Log
}

func NewBinlogReader(basepath string, time time.Time, idx int,
timeSpan int64, offset uint64, log *logger.Log) (binlog *BinlogReader, err error) {
    if 0 >= len(basepath) {
        err = fmt.Errorf("the argument is fail.")
        return
    }
    if !strings.HasSuffix(path, "/") {
        basepath = strings.Join([]string{basepath, string(os.PathSeparator)}, "")
    }
    binlog = new(BinlogReader)
    binlog.basepath = basepath
    binlog.idx = idx
    binlog.t = time
    binlog.timeSpan = timeSpan
    binlog.offset = offset
    binlog.log = log
    err = binlog.openBinlogReader()
    return
}

func (binlog *BinlogReader)openBinlogReader()(err error){
    for{
        filename := binlogFilename(binlog.t,binlog.basepath,binlog.idx)
        f,ok := os.Open(filename);
        if os.IsNotExist(ok){
            t := binlog.t.Add(time.Duration(binlog.timeSpan) * time.Second)
            if t.Before(time.Now()){
                binlog.t = t
                binlog.idx = 0;
                binlog.offset = 0
                continue
            } else{
                err = fmt.Errorf("the binlog is poll to current time,
                then no the binlog file.")
                return
            }
        } else {
            binlog.f = f;
            if(0 != binlog.offset){
                binlog.f.Seek(binlog.offset,os.SEEK_SET)
            }
            break
        }
    }
    return
}

func (binlog *BinlogReader) Close(isFlushStateFile bool) {
    if isFlushStateFile {
    }
    close(binlog.f)
}
func (binlog *Binlog) Read() (b []byte, err error) {
    l := int(unsafe.Sizeof(0)) //the int bytes length
    buf := make([]byte, BINLOG_READER_BUFFER_SIZE)
    for{
        var i int
        i,err = binlog.f.Read(buf)
        if nil != err{
            if(io.EOF != err) {
                return
            } else {
                binlog.f.Close(false)
                binlog.f.openBinlogReader()
                continue
            }
        }
        break;
    }
}
*/
