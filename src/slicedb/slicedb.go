package slicedb

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"time"
)
import "../skiplist"
import "../logger/binlog"
import "../fileio"

const (
	SlicedbStatusFileName      = "slicedb.status"
	ErrSliceIdxFileIdxOutScope = errors.New("slice idx file out of scope")
	SlicedbRuntimeFileName     = "slicedb.runtime"
)

type Slicepk struct { //the pk struct in the memory
	t      time.Time //the time of recode put into the system
	idx    int       // the file idx of the physical disk file
	offset int64
	len    int
	sum    []byte //the recode md5 sum
}
type SliceIdx struct {
	pk  *Slicepk
	sum []byte //the recode md5 sum
}
type SliceStatus struct {
	XMLName xml.Name  `xml:"slicedb"`
	begin   time.Time `xml:"begin"`
}

type SliceRuntime struct {
	scope time.Time `xml:"CurrentScope"`
	//the last complete slice pk and idx skiplist is flush to disk
	isflush bool `xml:"IsLastFlush"`
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
	datapath        string
	idxpath         string
	binlogpath      string
	statusfilepath  string
	runtimefilepath string

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
	if 0 == times {
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
