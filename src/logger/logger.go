package logger

import (
	"../octp_time"

	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	Debug = 0
	Info  = 1
	Warn  = 2
	Error = 3
	Fatal = 4
	Mark  = 100
)

type Logger struct {
	locker  sync.Mutex
	path    string
	level   int
	prefix  string
	maxsize uint64
	size    uint64
	f       *os.File
}

func New(path string, level int, prefix string, maxsize uint64) (log *Logger, err error) {
	if 0 == len(path) {
		err = fmt.Errorf("the argument:path is error.")
	}
	err = os.MkdirAll(path, 0777)
	if nil != err {
		if os.IsExist(err) {
			err = nil
		} else {
			return
		}
	}
	if !strings.HasSuffix(path, "/") {
		path = strings.Join([]string{path, string(os.PathSeparator)}, "")
	}
	log = new(Logger)
	log.path = path
	log.level = level
	log.prefix = prefix
	log.maxsize = maxsize
	err = log.open()
	return
}

func (log *Logger) open() (err error) {
	t := octp_time.GetNowShortString()
	p := strings.Join([]string{log.path, t, log.prefix, ".log"}, "")
	log.f, err = os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	return
}

func (log *Logger) Close() (err error) {
	if nil != log.f {
		err = log.f.Close()
	}
	log.size = 0
	return
}

func (log *Logger) Brush(level int, format string, args ...interface{}) (err error) {
	if level < log.level {
		return
	}
	line := fmt.Sprintf(format, args...)
	t := octp_time.GetNowLongString()
	var sLevel string
	switch level {
	case Debug:
		sLevel = "Debug!"
	case Info:
		sLevel = "Info!"
	case Warn:
		sLevel = "Warn!"
	case Error:
		sLevel = "Error!"
	case Fatal:
		sLevel = "Fatal!"
	case Mark:
		sLevel = "Mark!"
	}
	line = strings.Join([]string{sLevel, " ", t, " ", line, "\n"}, "")
	b := []byte(line)
	log.locker.Lock()
	defer log.locker.Unlock()
	if log.maxsize <= log.size {
		log.Close()
		log.open()
	}
	_, err = log.f.Write(b)
	if nil == err {
		log.size += uint64(len(b))
	}
	return
}
