package slicedb

import (
	"fmt"
	"os"
	"time"
)

func (db *Slicedb) SetDataPath(path string) (err error) {
	if 0 == len(path) {
		db.datapath = "/tmp/slicedb/data/"
		return
	}
	err = os.MkdirAll(path, os.ModePerm)
	if nil != err {
		if os.IsExist(err) {
			err = nil
		}
	}

	db.datapath = path
}

func (db *Slicedb) SetStatusFilePath(path string) (err error) {
	if 0 == len(path) {
		db.statusfilepath = "/tmp/slicedb/status/"
		return
	}
	if err = os.MkdirAll(path, os.ModePerm); nil != err {
		if os.IsExist(err) {
			err = nil
		}
	}
	db.statusfilepath = path
}
func (db *Slicedb) SetRuntimeFilePath(path string) (err error) {
	if 0 == len(path) {
		db.runtimefilepath = "/tmp/slicedb/status/"
		return
	}
	if err = os.MkdirAll(path, os.ModePerm); nil != err {
		if os.IsExist(err) {
			err = nil
		}
	}
	db.runtimefilepath = path
}

func (db *Slicedb) SetIdxPath(path string) (err error) {
	if 0 == len(path) {
		db.idxpath = "/tmp/sliedb/idx/"
		return
	}
	err = os.MkdirAll(path, os.ModePerm)
	if nil != err {
		if os.IsExist(err) {
			err = nil
		}
	}
	db.idxpath = path
}

func (db *Slicedb) SetBinlogPath(path string) (err error) {
	if 0 == len(path) {
		db.binlogpath = "/tmp/slicedb/binlog"
		return
	}
	err = os.MkdirAll(path, os.ModePerm)
	if nil != err {
		if os.IsExist(err) {
			err = nil
		}
	}
	db.binlogpath = path
}

func (db *Slicedb) SetPK(name string, idxtype int) (err error) {
	if 0 == len(name) {
		err = fmt.Errorf("the argument is null.")
		return
	}
	db.pkField = name
	db.pkType = idxtype
}

func (db *Slicedb) SetIdx0(name string, idxtype int) (err error) {
	if 0 == len(name) {
		err = fmt.Errorf("the argument is null.")
		return
	}
	db.idxField0 = name
	db.idxType0 = idxtype
}
func (db *Slicedb) SetIdx1(name string, idxtype int) (err error) {
	if 0 == len(name) {
		err = fmt.Errorf("the argument is null.")
		return
	}
	db.idxField1 = name
	db.idxType1 = idxtype
}
func (db *Slicedb) SetIdx2(name string, idxtype int) (err error) {
	if 0 == len(name) {
		err = fmt.Errorf("the argument is null.")
		return
	}
	db.idxField2 = name
	db.idxType2 = idxtype
}
func (db *Slicedb) SetSliceSizeInMemory(s int) {
	if 0 >= s {
		s = 30
	}
	db.maxSize = s
}
func (db *Slicedb) SetTimeSpan(s int64) {
	if 0 >= s {
		s = 300
	}
	db.ts = s
}
func (db *Slicedb) SetBufWriterSize(s int64) {
	if 0 >= s {
		s = 10 * 1024
	}
	db.bufsize = s
}
func (db *Slicedb) SetBufWriterTimeSpan(s int64) {
	if 0 >= s {
		s = 5
	}
	db.flush = s
}

func (db *Slicedb) setBeginDateTime(t time.Time) {
	if nil == t {
		t = GetStandardTime(time.Now())
	}
	db.begin = t
}
