/*
   the package is implements a lib for skiplist
   the skiplist is not thread-safe, if you want to
   thread-safe,then please make sure lock skiplist out the scope.
   the shiplist is not support duplicate values.
*/
package collections

import (
	"math/rand"
	"reflect"
	"time"
)

const (
	SKIPLISTNODENORMAL = 0
	SKIPLISTNODEHEAD   = 1
	SKIPLISTNODEFOOT   = 2
)

const SKIPLISTMAXLEVEL = 16

type SkipListNodeCmp interface {
	NodeKeyCmper(key1 interface{}, keylen1 uint32,
		key2 interface{}, keylen2 uint32) (rc int8)
}

type SkipList struct {
	Level    uint8
	MaxLevel uint8
	Head     []SkipListNode
	Foot     []SkipListNode
	Root     *SkipListNode
	Cmper    *SkipListNodeCmp
}

type SkipListNode struct {
	Level    uint64
	Property uint8
	Key      interface{}
	Keylen   uint32
	Value    interface{}
	Vallen   uint64
	Next     *SkipListNode
}

func SkipListCreate(maxLevel uint8, cmp *SkipListNodeCmp) (sl *SkipList, error err) {
	if 0 == maxLevel {
		sl = nil
		err = "create skiplist is fail."
		return
	}
	sl := new(SkipList)
	sl.MaxLevel = maxLevel
	sl.Level = 0
	sl.Root = nil
	sl.Head = make([]SkipListNode, maxLevel)
	sl.Foot = make([]SkipListNode, maxLevel)
	for i := 0; i < maxLevel; i++ {
		sl.Head[i].Property = SKIPLISTNODEHEAD
		sl.Foot[i].Property = SKIPLISTNODEFOOT
		sl.Head[i].Next = &(sl.Foot[i])
	}
	if nil == cmp {
		sl.Cmper = &(sl.NodeKeyCmp)
	} else {
		sl.Cmper = cmp
	}
	err = nil
	return
}

func randLevel() (level uint8) {
	r := rand.New(time.Now().UnixNano())
	level = r.Int63() % SKIPLISTMAXLEVEL
	return
}

func cmpResult(v int8) (rc int8) {
	switch rc {
	case rc > 0:
		return 1
	case rc == 0:
		return 0
	case rc < 0:
		return -1
	}
}
func intCmper(key1, key2 int64) int8 {
	return cmpResult(key1 - key2)
}
func uintCmper(key1, key2 uint64) int8 {
	return cmpResult(key1 - key2)
}
func floatCmper(key1, key2 float64) int8 {
	rc := key1 - key2
	switch rc {
	case rc > 0:
		return 1
	case rc == 0:
		return 0
	case rc < 0:
		return -1
	}
}

/*
   the string cmper is not suppering for Mars lanagreee
   and it only for normal word
*/
func stringCmper(key1 string, keylen1 uint32,
	key2 string, keylen2 uint32) int8 {
	if key1 > key2 {
		return 1
	} else if key1 < key2 {
		return -1
	} else {
		return 0
	}
	//please put into the code with rune
	//thanks huanshang
}

func timeCmper(key1, key2 string) uint8 {
}

func (sl *SkipList) NodeKeyCmp(key1 interface{}, keylen1 uint32,
	key2 interface{}, keylen2 uint32) (rc int8) {
	t1 := TypeOf(key1)
	t2 := TypeOf(key2)
	if t1 == t2 {
		switch {
		case t1 == reflect.Int:
			rc = key1 - key2
			return cmpResult(rc)
		}
	} else {
	}
}
func (sl *SkipList) Insert(key interface{}, uint32 keylen,
	val interface{}, uint64 vallen) {

}
