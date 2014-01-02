/*
   the package is implements a lib for skiplist
   the skiplist is not thread-safe, if you want to
   thread-safe,then please make sure lock skiplist out the scope.
   the shiplist is not support duplicate values.
*/
package skiplist

import (
	"../octp_time"

	"fmt"
	"math/rand"
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
	Level    int
	IdxType  int
	MaxLevel int
	Head     *SkipListNode
	Cmper    func(key1, key2 interface{}) (rc int8, err error)
}

type SkipListNode struct {
	Level    int
	Property int
	Key      interface{}
	Value    interface{}
	Next     []*SkipListNode
}

func randLevel() (level int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	level = int(r.Int63() % int64(SKIPLISTMAXLEVEL))
	return
}

func intCmper(key1, key2 interface{}) (rc int8, err error) {
	var k1, k2 int64
	var ok bool
	if k1, ok = key1.(int64); !ok {
		err = fmt.Errorf("the argument is not int type")
		return
	}
	if k2, ok = key2.(int64); !ok {
		err = fmt.Errorf("the argument is not int type")
		return
	}
	if k1 < k2 {
		rc = -1
		return
	} else if k1 > k2 {
		rc = 1
		return
	} else {
		rc = 0
		return
	}
	err = fmt.Errorf("the int key cannot compare")
	rc = 0
	return
}

func uintCmper(key1, key2 interface{}) (rc int8, err error) {
	var k1, k2 uint64
	var ok bool
	if k1, ok = key1.(uint64); !ok {
		err = fmt.Errorf("the argument is not uint type")
		return
	}
	if k2, ok = key2.(uint64); !ok {
		err = fmt.Errorf("the argument is not uint type")
		return
	}
	if k1 < k2 {
		rc = -1
		return
	} else if k1 > k2 {
		rc = 1
		return
	} else {
		rc = 0
		return
	}
	err = fmt.Errorf("the uint key cannot compare")
	rc = 0
	return
}

func floatCmper(key1, key2 interface{}) (rc int8, err error) {
	var k1, k2 float64
	var ok bool
	if k1, ok = key1.(float64); !ok {
		err = fmt.Errorf("the argument is not float type")
		return
	}
	if k2, ok = key2.(float64); !ok {
		err = fmt.Errorf("the argument is not float type")
		return
	}
	if k1 < k2 {
		rc = -1
		return
	} else if k1 > k2 {
		rc = 1
		return
	} else {
		rc = 0
		return
	}
	err = fmt.Errorf("the float key cannot compare")
	rc = 0
	return
}

/*
   the string cmper is not suppering for Mars lanagreee
   and it only for normal word
*/
func stringCmper(key1, key2 interface{}) (rc int8, err error) {
	var k1, k2 string
	var ok bool
	if k1, ok = key1.(string); !ok {
		err = fmt.Errorf("the argument is not string type")
		return
	}
	if k2, ok = key2.(string); !ok {
		err = fmt.Errorf("the argument is not string type")
		return
	}
	err = nil
	if k1 > k2 {
		rc = 1
		return
	} else if k1 < k2 {
		rc = -1
		return
	} else {
		rc = 0
		return
	}
	err = fmt.Errorf("the string key cannot compare")
	rc = 0
	return
	//please put into the code with rune
	//thanks huanshang
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
func SkipListCreate(maxLevel int, idxType int) (sl *SkipList, err error) {
	if 0 >= maxLevel {
		sl = nil
		err = fmt.Errorf("create skiplist is fail.")
		return
	}
	sl = new(SkipList)
	sl.MaxLevel = maxLevel
	sl.IdxType = idxType
	sl.Level = 0
	sl.Head = new(SkipListNode)
	sl.Head.Next = make([]*SkipListNode, maxLevel)
	sl.Head.Property = SKIPLISTNODEHEAD
	for i := 0; i < maxLevel; i++ {
		sl.Head.Next[i] = nil
	}
	switch idxType {
	case SKIPLIST_IDX_INT:
		sl.Cmper = intCmper
	case SKIPLIST_IDX_UINT:
		sl.Cmper = uintCmper
	case SKIPLIST_IDX_FLOAT:
		sl.Cmper = floatCmper
	case SKIPLIST_IDX_STRING:
		sl.Cmper = stringCmper
	case SKIPLIST_IDX_TIME:
		sl.Cmper = octp_time.TimeCmper
	default:
		sl.Cmper = stringCmper
	}
	err = nil
	return
}

func (sl *SkipList) Insert(key, val interface{}) (err error) {
	update := make([]*SkipListNode, sl.MaxLevel)
	var r int8
	p := sl.Head
	for i := sl.Level; i >= 0; i-- {
		for {
			q := p.Next[i]
			if nil == q {
				update[i] = p
				break
			} else {
				r, err = sl.Cmper(q.Key, key)
				if nil != err {
					return
				}
				if 0 > r {
					p = q
				} else if 0 == r {
					err = fmt.Errorf("the key is exist.")
					return
				} else {
					update[i] = p
					break
				}
			}
		}
	}

	k := randLevel()
	if k > sl.Level {
		for i := sl.Level + 1; i <= k; i++ {
			update[i] = sl.Head
		}
		sl.Level = k
	}
	n := new(SkipListNode)
	n.Key = key
	n.Level = k
	n.Value = val
	n.Next = make([]*SkipListNode, k+1)
	for i := 0; i <= k; i++ {
		n.Next[i] = update[i].Next[i]
		update[i].Next[i] = n
	}
	return
}

func (sl *SkipList) Find(key interface{}) (rc interface{}, err error) {
	var r int8
	p := sl.Head
	for i := sl.Level; i >= 0; i-- {
		for {
			q := p.Next[i]
			if nil == q {
				break
			} else {
				r, err = sl.Cmper(q.Key, key)
				if nil != err {
					return
				}
				if 0 > r {
					p = q
				} else if 0 == r {
					rc = q.Value
					return
				} else {
					break
				}
			}
		}
	}

	rc = nil
	err = fmt.Errorf("the object is not exist")
	return
}

func (sl *SkipList) Search(from, to interface{}) (
	rc map[interface{}]interface{}, err error) {
	update := make([]*SkipListNode, sl.MaxLevel)
	var r int8
	p := sl.Head
	for i := sl.Level; i >= 0; i-- {
		for {
			q := p.Next[i]
			if nil == q {
				update[i] = p
				break
			} else {
				r, err = sl.Cmper(q.Key, from)
				if nil != err {
					return
				}
				if 0 > r {
					p = q
				} else {
					update[i] = p
					break
				}
			}
		}
	}

	p = update[0]
	if nil == p {
		err = fmt.Errorf("not found keys slice")
		rc = nil
		return
	}
	var rc1, rc2 int8
	rc = make(map[interface{}]interface{})
	for {
		p = p.Next[0]
		if nil != p {
			rc1, err = sl.Cmper(p.Key, from)
			if nil != err {
				return
			}
			rc2, err = sl.Cmper(p.Key, to)
			if nil != err {
				return
			}
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

func (sl *SkipList) Print() (err error) {
	for i := sl.Level; i >= 0; i-- {
		p := sl.Head
		for {
			q := p.Next[i]
			if nil == q {
				break
			} else {
				switch sl.IdxType {
				case SKIPLIST_IDX_INT:
					var k, v int64
					var ok bool
					if k, ok = q.Key.(int64); !ok {
						return
					}
					if v, ok = q.Value.(int64); !ok {
						return
					}
					fmt.Printf("%d-%d  ", k, v)
				case SKIPLIST_IDX_UINT:
					fmt.Print("%10d(%10d)  ", q.Key, q.Value)
				case SKIPLIST_IDX_FLOAT:
					fmt.Print("%10f(%10f)  ", q.Key, q.Value)
				case SKIPLIST_IDX_TIME:
					fmt.Print("%10q(%10q)  ", q.Key, q.Value)
				case SKIPLIST_IDX_STRING:
					fmt.Print("%10q(%10q)  ", q.Key, q.Value)
				}
			}
			p = q
		}
		fmt.Println("")
	}
	return
}
