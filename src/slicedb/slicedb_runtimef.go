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
	fn := getFilepath(db.statusfilepath, SlicedbStatusFileName)
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
	fn := getFilepath(db.statusfilepath, SlicedbStatusFileName)
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

func (db *Slicedb) loadRuntimeFile() (err error) {
	if 0 == len(path) {
		err = fmt.Errorf("the path argument is empty.")
		return
	}
	if err = os.MkdirAll(path, 0777); nil != err {
		if os.IsExist(err) {
			err = nil
		}
	}
	fn := getFilepath(db.runtimefilepath, SlicedbRuntimeFileName)
	var c []byte
	if c, err = ioutil.ReadFile(fn); nil != err {
		return
	}
	var s SliceRuntime
	if err = xml.Unmarshal(c, &s); nil != err {
		return err
	}
}
