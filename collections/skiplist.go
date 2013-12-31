package collections

import (
	"fmt"
	"math/rand"
	"time"
)

const p = 0.25
const SKIPLISTMAXLEVEL = 15

//** Node Define **
type SkipListNode struct {
	Key, value interface{}
	next       []*SkipListNode
	previous   *SkipListNode
}

func (s *SkipListNode) Next() *SkipListNode {
	if len(s.next) <= 0 {
		return nil
	}

	return s.next[0]
}

func (s *SkipListNode) Previous() *SkipListNode {
	return s.previous
}

//** SkipList Define **
type SkipList struct {
	length   int
	MaxLevel int
	header   *SkipListNode
	footer   *SkipListNode
	compare  func(l, r interface{}) bool
}

func (s *SkipList) Length() int {
	return s.length
}

func (s *SkipList) Level() int {
	return len(s.header.next) - 1
}

//find the next node of key(or equal)
func (s *SkipList) find(key interface{}, forUpdate []*SkipListNode) *SkipListNode {
	current := s.header

	level := len(current.next) - 1
	if level < 0 {
		return nil
	}

	for i := level; i >= 0; i-- {
		for current.next[i] != nil && s.compare(current.next[i].Key, key) {
			current = current.next[i]
		}
		if forUpdate != nil {
			forUpdate[i] = current
		}
	}
	return current.Next()
}

func (s *SkipList) randomLevel() (n int) {
	for n = 0; n < s.Level()+1 && n < s.MaxLevel && rand.Float64() < p; n++ {
	}
	return
}

func (s *SkipList) Add(key, value interface{}) {
	currentLevel := s.Level()
	if key == nil {
		panic("key can not nil when insert")
	}
	forUpdate := make([]*SkipListNode, currentLevel+1, s.MaxLevel)
	findNode := s.find(key, forUpdate)

	if findNode != nil && findNode.Key == key {
		findNode.value = value
		return
	}

	newLevel := s.randomLevel()

	if newLevel > currentLevel {
		for i := currentLevel + 1; i <= newLevel; i++ {
			forUpdate = append(forUpdate, s.header)
			s.header.next = append(s.header.next, nil)
		}
	}

	newNode := &SkipListNode{
		Key:   key,
		value: value,
		next:  make([]*SkipListNode, newLevel+1, s.MaxLevel),
	}

	newNode.previous = forUpdate[0]

	for i := 0; i <= newLevel; i++ {
		newNode.next[i] = forUpdate[i].next[i]
		forUpdate[i].next[i] = newNode
	}

	s.length++

	if newNode.next[0] != nil {
		if newNode.next[0].previous != newNode {
			newNode.next[0].previous = newNode
		}
	}

	if s.footer == nil || s.compare(s.footer.Key, key) {
		s.footer = newNode
	}

	fmt.Println("Add key:", key, " ok. length:", s.length)

}

func (s *SkipList) Seek(key interface{}) *SkipListNode {
	if key == nil {
		return nil
	}

	findNode := s.find(key, nil)
	if findNode != nil && findNode.Key == key {
		return findNode
	}

	return nil
}

func (s *SkipList) Get(key interface{}) (value interface{}) {
	node := s.Seek(key)
	if node != nil {
		return node.value
	}

	return nil
}

func (s *SkipList) Remove(key interface{}) {
	currentLevel := s.Level()
	if key == nil {
		return
	}
	forUpdate := make([]*SkipListNode, currentLevel+1, s.MaxLevel)
	findNode := s.find(key, forUpdate)

	if findNode == nil || findNode.Key != key {
		return
	}

	if s.footer == findNode {
		s.footer = findNode.previous
	}

	findNode.previous.next[0] = findNode.next[0]
	if findNode.next[0] != nil {
		findNode.next[0].previous = findNode.previous
	}
	s.length--

	for i := len(forUpdate) - 1; i > 0; i-- {
		//fmt.Println("forUpdate:", i, forUpdate[i])
		//fmt.Println("findNode.next:", i, findNode.next[i])

		if forUpdate[i].next[i] == findNode {
			forUpdate[i].next[i] = findNode.next[i]
		}
	}

	for s.Level() > 0 && s.header.next[s.Level()] == nil {
		s.header.next = s.header.next[:s.Level()]
	}
	fmt.Println("Delete ", key, " ok. length:", s.length)
}

func New(compare func(l, r interface{}) bool) *SkipList {
	s := &SkipList{
		MaxLevel: SKIPLISTMAXLEVEL,
		header: &SkipListNode{
			next: []*SkipListNode{nil},
		},
		compare: compare,
	}

	return s
}

func NewInt() *SkipList {
	s := New(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})

	return s
}

func NewString() *SkipList {
	s := New(func(l, r interface{}) bool {
		return l.(string) < r.(string)
	})

	return s
}

func NewFloat64() *SkipList {
	s := New(func(l, r interface{}) bool {
		return l.(float64) < r.(float64)
	})

	return s
}
