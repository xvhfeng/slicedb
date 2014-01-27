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
	for i := 0; ; i++ {
		if r, err := binlog.NewReader(
			db.binlogPath, t, i, db.flushStatus.offset, log); nil != err {
			if nil != db.log {
				db.log.Brush(logger.Warn,
					"open the binlog file is fail,time:%s,idx:%d.",
					octp_time.GetDateTimeLongString(t), i)
				break
			}
		}
		if line, ok := r.Read(); nil != ok {
			if eos.IsEOF(ok) {
				return
			} else {
				if nil != db.log {
					db.log.Brush(logger.Error,
						"read binlog time:%s,idx:%d,err:%v.",
						octp_time.GetDateTimeLongString(t), i)
				}
			}
		}
		l := r.ExpLine(line)
		//    1:binlog OP|FILEIDX|OFFSET|LENGTH
	}
}

func (db *Slicedb) RestoreSliceByIdxFile(t time.Time) {

}
