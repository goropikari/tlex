package counter

type Counter struct {
	cnt int
}

func NewCounter() *Counter {
	return &Counter{cnt: 0}
}

func (c *Counter) Generate() int {
	c.cnt++
	return c.cnt
}
