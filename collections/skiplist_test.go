package collections

import (
	"fmt"
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
		fmt.Printf("%v: %v (level %d)\n", node.Key, node.value, len(node.next))
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

func TestInitialization(t *testing.T) {
	s := SkipList{
		MaxLevel: SKIPLISTMAXLEVEL,
		header: &SkipListNode{
			next: []*SkipListNode{nil},
		},
		compare: func(l, r interface{}) bool {
			return l.(int) < r.(int)
		},
	}

	s.Insert(1, "1")
	s.Insert(2, "2")
	s.Insert(3, "3")

	s.printRepr()
}
