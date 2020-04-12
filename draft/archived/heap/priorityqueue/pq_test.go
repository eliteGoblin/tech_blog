package priorityqueue

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	printPQ(pq)
	heap.Init(&pq)
	printPQ(pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	printPQ(pq)
	pq.update(item, item.value, 5)
	printPQ(pq)
	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}

func printPQ(pq PriorityQueue) {
	for _, v := range pq {
		fmt.Printf(" %+v ", *v)
	}
	fmt.Println("")
}
