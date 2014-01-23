package slicedb

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func getFilepath(path, file string) (filename string) {
	if strings.HasSuffix(path, "/") {
		filename = strings.Join([]string{path, string(os.PathSeparator)}, "")
	} else {
		filename = strings.Join([]string{path,
			string(os.PathSeparator), SlicedbStatusFileName}, "")
	}
	return filename
}

func (db *Slicedb) SaveStatus(path string, status *SliceStatus) error { /*{{{*/
	var err error
	if 0 == len(path) {
		return fmt.Errorf("the path argument is empty.")
	}

	if err = os.MkdirAll(path, os.ModePerm); nil != err {
		return err
	}
	s := &SliceStatus{begin: db.begin}
	var c []byte
	if c, err = xml.Marshal(s); nil != err {
		return err
	}
	fn := getFilepath(db.statusPath, SlicedbStatusFileName)
	if err = ioutil.WriteFile(fn, c, os.ModePerm); nil != err {
		return err
	}
	return nil
} /*}}}*/

func (db *Slicedb) loadStatus(path string) error { /*{{{*/
	var err error
	if 0 == len(path) {
		return fmt.Errorf("the path argument is empty.")
	}
	fn := getFilepath(db.statusPath, SlicedbStatusFileName)
	var c []byte
	if c, err = ioutil.ReadFile(fn); nil != err {
		return err
	}
	var s SliceStatus
	if err = xml.Unmarshal(c, &s); nil != err {
		return err
	}
	db.setBeginDateTime(s.begin)
	return nil
} /*}}}*/

func (db *Slicedb) loadFlushRuntime() error { /*{{{*/
	var err error
	if 0 == len(path) {
		return fmt.Errorf("the path argument is empty.")
	}
	fn := getFilepath(db.flushStatusPath, SlicedbFlushStatusFileName)
	var c []byte
	if c, err = ioutil.ReadFile(fn); nil != err {
		return err
	}
	if err = xml.Unmarshal(c, db.flushStatus); nil != err {
		return err
	}
} /*}}}*/

func (db *Slicedb) SaveFlushRuntime() error { /*{{{*/
	if 0 == len(path) {
		return fmt.Errorf("the path argument is empty.")
	}
	err = os.MkdirAll(path, 0777)
	if err = os.MkdirAll(db.flushStatusPath, os.ModePerm); nil != err {
		return err
	}
	var c []byte
	if c, err = xml.Marshal(&(db.flushStatus)); nil != err {
		return err
	}
	fn := getFilepath(db.flushStatusPath, SlicedbFlushStatusFileName)
	if err = ioutil.WriteFile(fn, c, os.ModePerm); nil != err {
		return err
	}
	return nil
} /*}}}*/
