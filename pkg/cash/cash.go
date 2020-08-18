package cash

// Cash - linked list
type Cash struct {
	top  *element
	back *element
}

// Push on back element
func (c *Cash) Push(data []int16) {
	element := &element{
		data: data,
	}

	if c.back != nil {
		c.back.next = element
	}

	c.back = element

	if c.top == nil {
		c.top = element
	}
}

// Pop return and delete element from top
func (c *Cash) Pop() (data []int16) {
	if c.top != nil {
		data = c.top.data
		c.top = c.top.next
	}
	return
}

// NewCash ...
func NewCash() *Cash {
	return &Cash{}
}
