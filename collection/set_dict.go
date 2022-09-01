package collection

import "crypto/sha256"

type Sha = [sha256.Size]byte

type BitsetDict[T any] struct {
	key  map[Sha]Bitset
	dict map[Sha]T
}

func NewBitsetDict[T any]() *BitsetDict[T] {
	return &BitsetDict[T]{
		key:  make(map[Sha]Bitset),
		dict: make(map[Sha]T),
	}
}

func (d *BitsetDict[T]) Set(bs Bitset, v T) *BitsetDict[T] {
	sha := sha256.Sum256(bs.x)
	d.key[sha] = bs
	d.dict[sha] = v

	return d
}

func (d *BitsetDict[T]) Get(bs Bitset) (T, bool) {
	sha := sha256.Sum256(bs.x)
	v, ok := d.dict[sha]
	return v, ok
}

// func (d *BitsetDict[T]) Keys() map[Sha]Bitset {
// 	return d.key
// }

// func (d *BitsetDict[T]) Dict() map[Sha]T {
// 	return d.dict
// }

func (d *BitsetDict[T]) Contains(bs Bitset) bool {
	sha := sha256.Sum256(bs.x)
	_, ok := d.dict[sha]
	return ok
}
