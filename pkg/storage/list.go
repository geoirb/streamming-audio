package storage

import (
	"container/list"
	"fmt"
	"io"
)

// Queue FIFO data struct
type queue struct {
	list *list.List
}

// Write on back element
func (q *queue) Write(data []byte) (n int, err error) {
	q.list.PushBack(data)
	return
}

// Read return and delete element from top
func (q *queue) Read(data []byte) (n int, err error) {
	element := q.list.Front()
	if element == nil {
		err = io.EOF
		return
	}
	fmt.Println(element)
	copy(data, element.Value.([]byte))
	q.list.Remove(element)
	return
}

func (q *queue) Close() (err error) {
	q.list = nil
	return
}
