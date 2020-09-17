package storage

import (
	"io"
)

type element struct {
	data []byte
	next *element
}

// Queue FIFO data struct
type queue struct {
	top  *element
	back *element
}

// Write on back element
// ATTETION!!! without copy
func (q *queue) Write(data []byte) (n int, err error) {
	element := &element{
		data: data,
	}

	if q.back != nil {
		q.back.next = element
	}

	q.back = element

	if q.top == nil {
		q.top = element
	}
	return
}

// Read return and delete element from top
func (q *queue) Read(data []byte) (n int, err error) {
	if q.top != nil {
		n = len(q.top.data)
		copy(data, q.top.data)
		q.top = q.top.next
		return n, nil
	}
	return 0, io.EOF
}

func (q *queue) Close() (err error) {
	q.top, q.back = nil, nil
	return
}
