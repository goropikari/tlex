package collection

// Algorithms, 4th Edition
// https://algs4.cs.princeton.edu/home/
type PriorityQueue[T any] struct {
	n        int
	data     []T
	lessFunc func(x, y T) bool
}

func NewPriorityQueue[T any](lessFunc func(x, y T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		n:        0,
		data:     make([]T, 1),
		lessFunc: lessFunc,
	}
}

func (pq *PriorityQueue[T]) Push(x T) {
	pq.data = append(pq.data, x)
	pq.n++
	pq.swim(pq.n)
}

func (pq *PriorityQueue[T]) Top() T {
	return pq.data[1]
}

func (pq *PriorityQueue[T]) Pop() {
	pq.swap(1, pq.n)
	pq.n--
	pq.data = pq.data[0 : pq.n+1]
	pq.sink(1)
}

func (pq *PriorityQueue[T]) Size() int {
	return pq.n
}

func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.n == 0
}

func (pq *PriorityQueue[T]) swim(k int) {
	for k > 1 && pq.less(k/2, k) {
		pq.swap(k/2, k)
		k /= 2
	}
}

func (pq *PriorityQueue[T]) sink(k int) {
	for 2*k <= pq.n {
		j := 2 * k
		if j < pq.n && pq.less(j, j+1) { // compare two children and select bigger one
			j++
		}
		if !pq.less(k, j) {
			break
		}
		pq.swap(k, j)
		k = j
	}
}

func (pq *PriorityQueue[T]) swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
}

func (pq *PriorityQueue[T]) less(i, j int) bool {
	return pq.lessFunc(pq.data[i], pq.data[j])
}
