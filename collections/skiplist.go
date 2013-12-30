package collections

import (
	"math/rand"
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
	for n = 0; n < s.Level()+1 && rand.Float64() < p; n++ {
	}
	return
}

func (s *SkipList) Insert(key, value interface{}) {
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

func (s *SkipList) Delete(key interface{}) {
	if key == nil {
		return
	}
	//forUpdate := make([]*SkipListNode, currentLevel+1, s.MaxLevel)
	//findNode = s.find(key, forUpdate)
}

type Comparable interface {
	Compare(l, r interface{}) bool
}
