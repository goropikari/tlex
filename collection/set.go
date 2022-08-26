package collection

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return Set[T]{}
}

func (s Set[T]) Insert(x T) Set[T] {
	s[x] = struct{}{}

	return s
}

func (s Set[T]) Erase(x T) Set[T] {
	delete(s, x)

	return s
}

func (s Set[T]) Contains(x T) bool {
	_, ok := s[x]
	return ok
}

func (s Set[T]) Difference(other Set[T]) Set[T] {
	d := NewSet[T]()
	for k := range s {
		if !other.Contains(k) {
			d.Insert(k)
		}
	}
	return d
}

func (s Set[T]) Intersection(other Set[T]) Set[T] {
	i := NewSet[T]()
	for k := range s {
		if other.Contains(k) {
			i.Insert(k)
		}
	}
	return i
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	u := NewSet[T]()
	for k := range s {
		u.Insert(k)
	}
	for k := range other {
		u.Insert(k)
	}
	return u
}

func (s Set[T]) Copy() Set[T] {
	t := NewSet[T]()
	for v := range s {
		t.Insert(v)
	}
	return t
}
