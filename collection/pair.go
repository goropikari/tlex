package collection

type Pair[T, W comparable] struct {
	First  T
	Second W
}

func NewPair[T, W comparable](first T, second W) Pair[T, W] {
	return Pair[T, W]{
		First:  first,
		Second: second,
	}
}
