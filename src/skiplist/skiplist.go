/*
   the package is implements a lib for skiplist
   the skiplist is not thread-safe, if you want to
   thread-safe,then please make sure lock skiplist out the scope.
   the shiplist is not support duplicate values.
*/
package skiplist

import (
	"math/rand"
	"slicedb/octp_time"
	"time"
)

const (
	SKIPLISTNODENORMAL = 0
	SKIPLISTNODEHEAD   = 1
	SKIPLISTNODEFOOT   = 2

	SKIPLIST_IDX_INT    = 0
	SKIPLIST_IDX_UINT   = 1
	SKIPLIST_IDX_STRING = 2
	SKIPLIST_IDX_FLOAT  = 4
	SKIPLIST_IDX_TIME   = 8
)

const SKIPLISTMAXLEVEL = 16

type SkipList struct {
	Level    uint8
	IdxType  uint8
	MaxLevel uint8
	Head     []SkipListNode
	/* Foot     []SkipListNode */
	/* Root     *SkipListNode */
	Cmper func(key1, key2 interface{}) uint8
}

type SkipListNode struct {
	Level    uint64
	Property uint8
	Key      interface{}
	Keylen   uint32
	Value    interface{}
	Vallen   uint64
	Next     []SkipListNode
}

/*
huanshan:
    the func name is XXXCreate not XXXNew,because from c.
    in the c,create mapping destroy,init mapping free.
    because create means create a memory ptr and init the memory,
    destroy means free the memory ptr and set null to ptr;
    init means it exist a memory ptr(in stack or heap)
    then set 0 to the memory;
    free means release the memory and the ptr is not change.
    new?? there is no New func in C
*/
func SkipListCreate(maxLevel uint8, idxType uint8) (sl *SkipList, error err) {
	if 0 == maxLevel {
		sl = nil
		err = "create skiplist is fail."
		return
	}
	sl := new(SkipList)
	sl.MaxLevel = maxLevel
	sl.IdxType = idxType
	sl.Level = 0
	sl.Root = nil
	sl.Head = make([]SkipListNode, maxLevel)
	sl.Foot = make([]SkipListNode, maxLevel)
	for i := 0; i < maxLevel; i++ {
		sl.Head[i].Property = SKIPLISTNODEHEAD
		/* sl.Foot[i].Property = SKIPLISTNODEFOOT */
		sl.Head[i].Next = SKIPLISTFOOT
		/* sl.Foot[i].Next = nil */
	}
	switch idxType {
	case SKIPLIST_IDX_INT:
		sl.Cmper = intCmper
	case SKIPLIST_IDX_UINT:
		sl.Cmper = uintCmper
	case SKLIPLIST_IDX_FLOAT:
		sl.Cmper = floatCmper
	case SKIPLIST_IDX_STRING:
		sl.Cmper = stringCmper
	case SKIPLIST_IDX_TIME:
		sl.Cmper = octp_time.timeCmper
	default:
		sl.Cmper = stringCmper
	}
	err = nil
	return
}

func randLevel() (level uint8) {
	r := rand.New(time.Now().UnixNano())
	level = r.Int63() % SKIPLISTMAXLEVEL
	return
}

func intCmper(key1, key2 int64) int8 {
	if key1 < key2 {
		return -1
	} else if key1 > key2 {
		return 1
	} else {
		return 0
	}
}

func uintCmper(key1, key2 uint64) int8 {
	if key1 < key2 {
		return -1
	} else if key1 > key2 {
		return 1
	} else {
		return 0
	}
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

func (sl *SkipList) Insert(key interface{}, uint32 keylen,
	val interface{}, uint64 vallen) (rc bool, err error) {
	update := make(*SkipListNode, sl.MaxLevel)
	p := &(sl.Head[sl.Level])
	for i := sl.Level; i >= 0; i-- {
		for {
			q := &(p.Next[i])
			if nil == q {
				break
			} else {
				r := sl.Cmper(q.Key, key)
				if 0 > r {
					ipdate[i] = q
				} else if 0 == r {
					err = "the key is exist."
					rc = false
					return
				} else {
					break
				}
			}
			p = q
		}
	}

	k := randLevel()
	if k > sl.Level {
		for i := sl.Level; i < k; i++ {
			update[i] = &(sl.Header[i])
		}
		sl.Level = k
	}

	n := new(SkipListNode)
	n.Key = key
	n.Level = k
	n.Value = val
	n.Next = make(*SkipListNode, k)

	for i := k; i >= 0; i-- {
		p := update[i]
		n.Next[i] = p.Next[i]
		p.Next[i] = n
	}
	rc = true
	return
}

func (sl *SkipList) Find(key interface{}) (rc interface{}, err error) {
	update := make(*SkipListNode, sl.MaxLevel)
	p := &(sl.Head[sl.Level])
	for i := sl.Level; i >= 0; i-- {
		for {
			q := &(p.Next[i])
			if nil == q {
				break
			} else {
				r := sl.Cmper(q.Key, from)
				if 0 == r {
					rc = q.Value
					return
				} else if 0 < r {
					break
				}
			}
			p = q
		}
	}
	rc = nil
	err = "the object is not exist"
	return
}

func (sl *SkipList) Search(from, to interface{}) (
	rc map[interface{}]interface{}, err error) {

	update := make(*SkipListNode, sl.MaxLevel)
	p := &(sl.Head[sl.Level])
	for i := sl.Level; i >= 0; i-- {
		for {
			q := &(p.Next[i])
			if nil == q {
				break
			} else {
				r := sl.Cmper(q.Key, from)
				if 0 > r {
					update[i] = q
				} else {
					break
				}
			}
			p = q
		}
	}

	p = update[0]
	if nil == p {
		err = "not found keys slice"
		rc = nil
		return
	}
	for {
		p = p.Next[0]
		if nil != p {
			rc1 := sl.Cmper(p.Key, from)
			rc2 := sl.Cmper(p.Key, to)
			if 0 <= rc1 && 0 >= rc2 {
				rc[p.Key] = p.Value
			} else {
				break
			}
		} else {
			break
		}
	}
	err = nil
	return
}
