package collection

// https://github.com/atcoder/live_library/blob/2f63329d476ccb6dfe3d60e2fde468b8c8797dda/uf.cpp

type UnionFind struct {
	d []int
}

func NewUnionFind(n int) *UnionFind {
	d := make([]int, n)
	for i := 0; i < n; i++ {
		d[i] = -1
	}
	return &UnionFind{d: d}
}

func (uf *UnionFind) Find(x int) int {
	if uf.d[x] < 0 {
		return x
	}
	uf.d[x] = uf.Find(uf.d[x])
	return uf.d[x]
}

func (uf *UnionFind) Unite(x, y int) bool {
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return false
	}
	if uf.d[x] > uf.d[y] {
		x, y = y, x
	}
	uf.d[x] += uf.d[y]
	uf.d[y] = x
	return true
}

func (uf *UnionFind) Same(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFind) Size(x int) int {
	return -uf.d[uf.Find(x)]
}

func (uf *UnionFind) Group() [][]int {
	memo := make(map[int][]int)
	for i := range uf.d {
		memo[uf.Find(i)] = append(memo[uf.Find(i)], i)
	}
	ret := make([][]int, 0, len(memo))
	for _, arr := range memo {
		ret = append(ret, arr)
	}
	return ret
}
