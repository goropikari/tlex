package collection

type Tuple[T, W comparable] struct {
	First  T
	Second W
}

func NewTuple[T, W comparable](first T, second W) Tuple[T, W] {
	return Tuple[T, W]{
		First:  first,
		Second: second,
	}
}
