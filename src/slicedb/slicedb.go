package slicedb

import (
	"fmt"
	"os"
	"time"
)
import "../skiplist"
import "../logger/binlog"
import "../fileio"

type Slice struct {
	pk   *skiplist.SkipList
	idx0 *skiplist.SkipList
	idx1 *skiplist.SkipList
	idx2 *skiplist.SkipList

	size int64 // the slice-chunk file size//the slice timespan
	idx  int   //the slice-chunk file idx
	t    time.Time
	end  time.Time // the slice end time

	f      *os.File
	buf    *fileio.BufWriter
	binlog *binlog.Binlog
	log    *logger.Log
}

type Slicedb struct {
	datapath   string
	idxpath    string
	binlogpath string

	pkField   string
	idxField0 string
	idxField1 string
	idxField2 string
	pkType    int
	idxType0  int
	idxType1  int
	idxType2  int

	maxSize int //the slices idx count in the memory
	size    int //the current slices idx count in the memory
	idxs    *skiplist.ShipList
	ts      int //the slice timespan
	// if in c,i store size fd with mmap for reads
	bufsize int64 //the bufriter buffer size

	begin time.Time
	flush int //the file flush into io timespan
	log   *logger.Log
}

func NewSlicedb(log *logger.Log) (db *Slicedb, err error) {
	db = new(Slicedb)
	db.log = log
}

func (db *Slicedb) SetDataPath(path string) (err error) {
	if 0 == len(path) {
		db.datapath = "/tmp/slicedb/data/"
		return
	}
	err = os.MkdirAll(path, 0777)
	if nil != err {
		if os.IsExist(err) {
			err = nil
		}
	}

	db.datapath = path
}

func (db *Slicedb) SetIdxPath(path string) (err error) {
	if 0 == len(path) {
		db.idxpath = "/tmp/sliedb/idx/"
		return
	}
	err = os.MkdirAll(path, 0777)
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
	err = os.MkdirAll(path, 0777)
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

/*
func (db *Slicedb)RunTimes(begin,now time.Time)(m,r int64){
	if begin.After(now) {
		err = fmt.Errorf("the begin time before now.")
		return
	}
	secs := int64(now.Sub(begin).Seconds())
	m = secs / ts
    r = secs % ts
}
*/
func (db *Slicedb) BaseTime(begin time.Time, now time.Time,
	ts int64) (t time.Time, err error) {
	if begin.After(now) {
		err = fmt.Errorf("the begin time before now.")
		return
	}
	secs := int64(now.Sub(begin).Seconds())
	m := secs / ts
	t = begin.Add(time.Duration(m * ts * time.Second))
	return
}

func (db *Slicedb) SetBeginDateTime(t time.Time) {
	if nil == t {
		t = time.Now()
	}
	db.begin = t
}

func (db *Slicedb) Start() (err error) {
	//restore the data into memory
	n := time.Now()
	t, ok := db.BaseTime(db.begin, n, db.ts)

	//open the current db
}
