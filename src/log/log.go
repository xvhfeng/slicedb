package logger

import (
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
	File    *f
}

func New(path string, perm FileMode, level int, prefix string, maxsize uint64) (log *logger, error err) {
	if nil == path || 0 == len(path) {
		err = fmt.Errorf("the argument:path is error.")
	}
	err = os.MkdirAll(path, perm)
	if nil != err {
		if os.IsExist(err) {
			err = nil
		} else {
			return
		}
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	log = new(Logger)
	log.path = path
	log.level = level
	log.prefix = prefix
	log.maxsize = maxisze
	err = log.open()
	return
}

func (log *Logger) open() (err error) {
	t := time.GetNowShortString()
	p := log.path + t + log.prefix + ".log"
	log.f, err = os.Open(p, os.O_CREATE|os.O_APPEND)
	if nil != err {
		return
	}
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
	line := fmt.Sprintf(format, args)
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

	line = t + " " + sLevel + line + "\n"
	log.locker.Lock()
	defer log.locker.Unlock()
	if log.maxsize <= log.size {
		log.Close()
		log.open()
	}
	_, err = log.f.WriteString(line)
	if nil == err {
		log.size += len(line)
	}
	return
}
