package storage

type element struct {
	data []byte
	next *element
}

// List data struct
type list struct {
	top  *element
	back *element
}

// Push on back element
func (l *list) Push(data []byte) {
	element := &element{
		data: make([]byte, len(data)),
	}
	copy(element.data, data)

	if l.back != nil {
		l.back.next = element
	}

	l.back = element

	if l.top == nil {
		l.top = element
	}
}

// Pop return and delete element from top
func (l *list) Pop() (data []byte) {
	if l.top != nil {
		data = l.top.data
		l.top = l.top.next
	}

	return
}
