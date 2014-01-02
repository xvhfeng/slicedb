package skiplist

import (
	"fmt"
	"syscall"
	"testing"
)

func TestSkipListCreate(t *testing.T) {

	maxLevel := 16
	var idxType int = SKIPLIST_IDX_INT
	sl, err := SkipListCreate(maxLevel, idxType)
	if nil != err {
		t.Error(err)
		return
	}
	/* fmt.Println("idx type:%d,max-level:%d,level:%d.", */
	/* sl.IdxType, sl.MaxLevel, sl.Level) */
	var k1, v1 int64
	var k2, v2 int64
	k1, v1 = 1, 1
	k2, v2 = 10, 10
	sl.Insert(k1, v1)
	/* sl.Print() */
	/* fmt.Printf("Head") */
	/* fmt.Println(sl.Head) */
	err1 := sl.Insert(k2, v2)
	if nil != err1 {
		fmt.Println(err1)
	}
	var k3, v3 int64
	k3, v3 = 100, 100
	sl.Insert(k3, v3)
	var k4, v4 int64
	k4, v4 = 20, 20
	sl.Insert(k4, v4)
	var k5, v5 int64
	k5, v5 = 40, 40
	sl.Insert(k5, v5)
	var k6, v6 int64
	k6, v6 = 30, 30
	sl.Insert(k6, v6)
	sl.Print()
	var v interface{}
	var k int64
	k = 40
	v, _ = sl.Find(k)
	fmt.Println(v)
	var kf, kt int64
	kf, kt = 10, 50

	/* var rc map[interface{}]interface{} */
	rc, _ := sl.Search(kf, kt)
	for _, vv := range rc {
		fmt.Println(vv)
	}
	e := syscall.Mkdir("/opt/gogogogog/path/ath", 777)
	if nil != e {
		fmt.Println(e)
	}
}
