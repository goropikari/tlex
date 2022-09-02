package counter

type Counter struct {
	cnt int
}

func NewCounter(start int) *Counter {
	return &Counter{cnt: start - 1}
}

func (c *Counter) Generate() int {
	c.cnt++
	return c.cnt
}
