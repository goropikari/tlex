package collection

import "errors"

type Bitset struct {
	length int
	x      []byte
}

func NewBitset(n int) Bitset {
	return Bitset{
		length: n,
		x:      make([]byte, (n+7)/8),
	}
}

func (b Bitset) Bytes() []byte {
	return b.x
}

func (b Bitset) GetLength() int {
	return b.length
}

func (b Bitset) Copy() Bitset {
	x := make([]byte, (b.length+7)/8)
	copy(x, b.x)

	return Bitset{
		length: b.length,
		x:      x,
	}
}

func (b Bitset) Up(n int) Bitset {
	id := n / 8
	pos := n % 8
	nb := b.Copy()
	nb.x[id] |= (1 << pos)

	return nb
}

func (b Bitset) Union(o Bitset) Bitset {
	if b.length != o.length {
		panic(errors.New("different bitset size"))
	}
	b = b.Copy()
	for i, v := range o.x {
		b.x[i] |= v
	}
	return b
}

func (b Bitset) Intersection(o Bitset) Bitset {
	if b.length != o.length {
		panic(errors.New("different bitset size"))
	}
	b = b.Copy()
	for i, v := range o.x {
		b.x[i] &= v
	}
	return b
}

func (b Bitset) Contains(n int) bool {
	id := n / 8
	pos := n % 8
	return (b.x[id] & (1 << pos)) > 0
}

func (b Bitset) IsZero() bool {
	for _, v := range b.x {
		if v > 0 {
			return false
		}
	}
	return true
}
