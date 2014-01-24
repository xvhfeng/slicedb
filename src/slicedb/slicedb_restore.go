package slicedb

import "time"

func (db *Slicedb) RestoreSlice(t time.Time) {
	for i := 0; ; i++ {
		if t.Before(db.flushStatus.scope) { // read binlog

		} else { //read idx-file
		}
	}
}

func (db *Slicedb) RestoreSliceByBinlog(t time.Time) {

}

func (db *Slicedb) RestoreSliceByIdxFile(t time.Time) {

}
