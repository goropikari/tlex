package collection

var noelems = struct{}{}

type Set[T comparable] struct {
	Mp    map[T]struct{}
	Elems []T
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		Mp:    make(map[T]struct{}),
		Elems: make([]T, 0),
	}
}

func (s *Set[T]) Size() int {
	return len(s.Elems)
}

func (s *Set[T]) Insert(x T) *Set[T] {
	_, ok := s.Mp[x]
	if !ok {
		s.Mp[x] = struct{}{}
		s.Elems = append(s.Elems, x)
	}

	return s
}

func (s *Set[T]) Erase(x T) *Set[T] {
	delete(s.Mp, x)
	elems := make([]T, 0, len(s.Elems))
	for _, v := range s.Elems {
		if v == x {
			continue
		}
		elems = append(elems, v)
	}
	s.Elems = elems

	return s
}

func (s *Set[T]) Contains(x T) bool {
	_, ok := s.Mp[x]
	return ok
}

func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	d := NewSet[T]()
	for _, k := range s.Elems {
		if !other.Contains(k) {
			d.Insert(k)
		}
	}
	return d
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	i := NewSet[T]()
	for _, k := range s.Elems {
		if other.Contains(k) {
			i.Insert(k)
		}
	}
	return i
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	u := NewSet[T]()
	for _, k := range s.Elems {
		u.Insert(k)
	}
	for _, k := range other.Elems {
		u.Insert(k)
	}
	return u
}

func (s *Set[T]) Copy() *Set[T] {
	t := NewSet[T]()
	for _, v := range s.Elems {
		t.Insert(v)
	}
	return t
}

func (s *Set[T]) Slice() []T {
	return s.Elems
}

func (s *Set[T]) Iterator() *setIterator[T] {
	return &setIterator[T]{
		currIdx: 0,
		length:  len(s.Elems),
		elems:   s.Elems,
	}
}

type setIterator[T comparable] struct {
	currIdx int
	length  int
	elems   []T
}

func (iter *setIterator[T]) HasNext() bool {
	return iter.currIdx < iter.length
}

func (iter *setIterator[T]) Next() T {
	ret := iter.elems[iter.currIdx]
	iter.currIdx++
	return ret
}
