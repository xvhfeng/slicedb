package slicedb

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)
import "../skiplist"
import "../logger/binlog"
import "../fileio"

const (
	SlicedbStatusFileName      = "slicedb.status"
	ErrSliceIdxFileIdxOutScope = errors.New("slice idx file out of scope")
)

type SliceStatus struct {
	XMLName xml.Name  `xml:"slicedb"`
	begin   time.Time `xml:"begin"`
}

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
	datapath       string
	idxpath        string
	binlogpath     string
	statusfilepath string

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

func (db *Slicedb) RunTimeAndCurrentScope(begin, now time.Time) (m int64, t time.Time) {
	if begin.After(now) {
		err = fmt.Errorf("the begin time before now.")
		return
	}
	secs := int64(now.Sub(begin).Seconds())
	m = secs / db.ts
	t = GetStandardTime(begin).Add(time.Duration(m * time.Second * db.ts))
	return
}
func (db *Slicedb) BaseTime(begin time.Time, now time.Time,
	ts int64) (t time.Time, err error) {
	if begin.After(now) {
		err = fmt.Errorf("the begin time before now.")
		return
	}
	secs := int64(now.Sub(begin).Seconds())
	m := secs / ts
	t = GetStandardTime(begin).Add(time.Duration(m * ts * time.Second))
	return
}

func (db *Slicedb) setBeginDateTime(t time.Time) {
	if nil == t {
		t = GetStandardTime(time.Now())
	}
	db.begin = t
}

func (db *Slicedb) SaveStatusFile(path string, status *SliceStatus) (err error) { /*{{{*/
	if 0 == len(path) {
		err = fmt.Errorf("the path argument is empty.")
		return
	}
	err = os.MkdirAll(path, 0777)
	if nil != err {
		if os.IsExist(err) {
			err = nil
		}
	}
	s := &SliceStatus{begin: db.begin}
	var c []byte
	if c, err = xml.Marshal(s); nil != err {
		return err
	}

	var fn string
	if strings.HasSuffix(path, "/") {
		fn = strings.Join([]string{path, string(os.PathSeparator)}, "")
	} else {
		fn = strings.Join([]string{path,
			string(os.PathSeparator), SlicedbStatusFileName}, "")
	}

	if err = ioutil.WriteFile(fn, c, os.ModePerm); nil != err {
		return
	}
} /*}}}*/

func (db *Slicedb) loadStatusFile(path string) (err error) { /*{{{*/
	if 0 == len(path) {
		err = fmt.Errorf("the path argument is empty.")
		return
	}
	err = os.MkdirAll(path, 0777)
	if nil != err {
		if os.IsExist(err) {
			err = nil
		}
	}
	var fn string
	if strings.HasSuffix(path, "/") {
		fn = strings.Join([]string{path, string(os.PathSeparator)}, "")
	} else {
		fn = strings.Join([]string{path,
			string(os.PathSeparator), SlicedbStatusFileName}, "")
	}
	var c []byte
	if c, err = ioutil.ReadFile(fn); nil != err {
		return
	}
	var s SliceStatus
	if err = xml.Unmarshal(c, &s); nil != err {
		return err
	}
	db.setBeginDateTime(s.begin)
	return
} /*}}}*/

func GetStandardTime(t time.Time) (st time.Time) {
	st = time.Date(t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), 0, time.UTC)
	return
}

func SliceNow() (t time.Time) {
	t = GetStandardTime(time.Now())
	return
}

func (db *Slicedb) NewSlice() (err error) {

}
func (db *Slicedb) RestoreSlice(t time.Time) {
	for i := 0; ; i++ {

	}
}
func (db *Slicedb) Start() (err error) {
	/* 1:compute the times of putting data into memory form begin start time
		    2:compute the last basetime to now
		    3;judge the current now is in the last basetime scope
	        4:if 3 is true,load the current data from data-file into index
	        5:if 3 is false,load the last idx-file into memory
	*/
	if db.idxs, err =
		skiplist.New(skiplist.SKIPLISTMAXLEVEL,
			skiplist.SKIPLIST_IDX_TIME); nil != err {
		return
	}
	// 1
	now := GetStandardTime(time.Now())
	if err = db.loadStatusFile(db.statusfilepath); nil != err {
		return err
	} else {
		db.setBeginDateTime(now)
	}
	times, scope := db.RunTimesAndCurrentScope(db.begin, now)
	if 0 == times { //the slicedb is the first starting
		err = db.RestoreSlice(scope)
		err = db.NewSlice()
		return err
	}

	var t int
	if times > db.maxSize {
		t = db.maxSize
	} else {
		t = times
	}
	for t > 0 {

	}
}
