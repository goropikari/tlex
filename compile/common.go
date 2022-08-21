package compile

type Tuple[T, W comparable] struct {
	first  T
	second W
}

func NewTuple[T, W comparable](first T, second W) Tuple[T, W] {
	return Tuple[T, W]{
		first:  first,
		second: second,
	}
}
