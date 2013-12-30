package collections

import (
	"fmt"
	"strconv"
	"testing"
)

func (s *SkipList) printRepr() {

	fmt.Printf("header:\n")
	for i, link := range s.header.next {
		if link != nil {
			fmt.Printf("\t%d: -> %v\n", i, link.Key)
		} else {
			fmt.Printf("\t%d: -> END\n", i)
		}
	}

	for node := s.header.Next(); node != nil; node = node.Next() {
		fmt.Printf("%v: %v (level %d)\n", node.Key, node.value, len(node.next)-1)
		for i, link := range node.next {
			if link != nil {
				fmt.Printf("\t%d: -> %v\n", i, link.Key)
			} else {
				fmt.Printf("\t%d: -> END\n", i)
			}
		}
	}
	fmt.Println()
}

func (s *SkipList) printList() {
	for i := s.Level(); i >= 0; i-- {
		fmt.Printf("level%d:", i)
		for node := s.header; node != nil; node = node.next[i] {
			if node.next[i] != nil {
				fmt.Printf("%v->", node.next[i].Key)
			} else {
				fmt.Println("END")
			}
		}
	}
}

func TestInsert(t *testing.T) {
	s := initIntList()

	s.printList()
}

func initIntList() SkipList {
	s := SkipList{
		MaxLevel: SKIPLISTMAXLEVEL,
		header: &SkipListNode{
			next: []*SkipListNode{nil},
		},
		compare: func(l, r interface{}) bool {
			return l.(int) < r.(int)
		},
	}

	for i := 1; i <= 1000; i++ {
		s.Insert(i, "value"+strconv.Itoa(i))
	}

	return s
}

func TestSeek(t *testing.T) {
	s := initIntList()
	i := 381
	node := s.Seek(i)
	fmt.Println(strconv.Itoa(i), ":", node)
}
