package automata

import (
	"crypto/sha256"
	stdmath "math"
	"unicode"

	"github.com/goropikari/tlex/collection"
	"github.com/goropikari/tlex/utils/guid"
)

const asciiSize = 1 << 7

const epsilon = byte(255)
const nonFinStateRegexID RegexID = stdmath.MaxInt

var SupportedChars = []byte{}

func init() {
	for i := 1; i < asciiSize; i++ {
		SupportedChars = append(SupportedChars, byte(i))
	}
}

var UnicodeRange = []Interval{
	NewInterval(0, int(unicode.MaxRune)),
}

type RegexID int
type StateID int
type Sha = [sha256.Size]byte
type Nothing struct{}

var nothing = Nothing{}

type Interval struct {
	L int
	R int
}

func NewInterval(s, e int) Interval {
	return Interval{
		L: s,
		R: e,
	}
}

func (x Interval) Overlap(y Interval) bool {
	return y.L <= x.R && x.L <= y.R
}

func (x Interval) Difference(y Interval) []Interval {
	if !x.Overlap(y) {
		return []Interval{x}
	}

	ret := make([]Interval, 0, 2)
	if x.L < y.L {
		ret = append(ret, NewInterval(x.L, y.L-1))
	}
	if x.R > y.R {
		ret = append(ret, NewInterval(y.R+1, x.R))
	}

	return ret
}

// https://stackoverflow.com/a/25832898
func Disjoin(intvs []Interval) []Interval {
	pq := collection.NewPriorityQueue(func(x, y Interval) bool {
		// ascending order
		if x.L != y.L {
			return x.L > y.L
		}
		return x.R > y.R
	})

	for _, v := range intvs {
		pq.Push(v)
	}

	ret := make([]Interval, 0, len(intvs))
	for pq.Size() >= 2 {
		t1 := pq.Top()
		pq.Pop()
		t2 := pq.Top()
		pq.Pop()

		if t1.Overlap(t2) {
			if t1.L < t2.L {
				nx1 := NewInterval(t1.L, t2.L-1)
				nx2 := NewInterval(t2.L, t1.R)
				nx3 := NewInterval(t2.L, t2.R)
				pq.Push(nx1)
				pq.Push(nx2)
				pq.Push(nx3)
			} else { // t1.l == t2.l
				pq.Push(t1)
				nx := NewInterval(t1.R+1, t2.R)
				if t1.R+1 <= t2.R {
					pq.Push(nx)
				}
			}
		} else {
			ret = append(ret, t1)
			pq.Push(t2)
		}
	}
	ret = append(ret, pq.Top())

	return ret
}

func NewStateID() StateID {
	return StateID(guid.New())
}

type StateIDToRegexID map[StateID]RegexID

func NewStateIDToRegexID() StateIDToRegexID {
	return make(StateIDToRegexID)
}

func (mp StateIDToRegexID) Get(sid StateID) RegexID {
	v, ok := mp[sid]
	if ok {
		return v
	}

	return nonFinStateRegexID
}

func (mp StateIDToRegexID) Set(sid StateID, rid RegexID) {
	mp[sid] = rid
}

type State struct {
	id StateID
}

func NewState(id StateID) State {
	return State{id: id}
}

func (st State) GetID() StateID {
	return st.id
}

type StateSet struct {
	bs collection.Bitset
}

func NewStateSet(n int) *StateSet {
	return &StateSet{
		bs: collection.NewBitset(n),
	}
}

func (ss *StateSet) Insert(x StateID) *StateSet {
	ss.bs = ss.bs.Up(int(x))
	return ss
}

func (ss *StateSet) Copy() *StateSet {
	return &StateSet{
		bs: ss.bs.Copy(),
	}
}

func (ss *StateSet) Union(other *StateSet) *StateSet {
	return &StateSet{
		bs: ss.bs.Union(other.bs),
	}
}

func (ss *StateSet) Intersection(other *StateSet) *StateSet {
	return &StateSet{
		bs: ss.bs.Intersection(other.bs),
	}
}

func (ss *StateSet) IsAny() bool {
	return !ss.IsEmpty()
}

func (ss *StateSet) IsEmpty() bool {
	return ss.bs.IsZero()
}

func (ss *StateSet) Contains(x StateID) bool {
	return ss.bs.Contains(int(x))
}

type stateSetIterator struct {
	maxID  StateID
	currID StateID
	ss     *StateSet
}

func newStateSetIterator(ss *StateSet) *stateSetIterator {
	sid := StateID(0)
	maxID := StateID(ss.bs.GetLength() - 1)
	for sid <= maxID {
		if ss.Contains(sid) {
			break
		}
		sid++
	}

	return &stateSetIterator{
		maxID:  maxID,
		currID: sid,
		ss:     ss,
	}
}

func (iter *stateSetIterator) HasNext() bool {
	return iter.currID <= iter.maxID
}

func (iter *stateSetIterator) Next() StateID {
	ret := iter.currID
	iter.currID++
	for iter.currID <= StateID(iter.maxID) {
		if iter.ss.Contains(iter.currID) {
			break
		}
		iter.currID++
	}

	return ret
}

func (ss *StateSet) iterator() *stateSetIterator {
	if ss == nil {
		return &stateSetIterator{
			maxID:  -1,
			currID: 0,
			ss:     nil,
		}
	}

	return newStateSetIterator(ss)
}

type StateSetDict[T any] struct {
	dict    *collection.BitsetDict[T]
	shaToSs map[Sha]*StateSet
}

func NewStateSetDict[T any]() *StateSetDict[T] {
	return &StateSetDict[T]{
		dict:    collection.NewBitsetDict[T](),
		shaToSs: make(map[Sha]*StateSet),
	}
}

func (d *StateSetDict[T]) Set(ss *StateSet, v T) *StateSetDict[T] {
	d.dict.Set(ss.bs, v)
	sha := sha256.Sum256(ss.bs.Bytes())
	d.shaToSs[sha] = ss
	return d
}

func (d *StateSetDict[T]) Get(ss *StateSet) (T, bool) {
	v, ok := d.dict.Get(ss.bs)
	return v, ok
}

func (d *StateSetDict[T]) Contains(ss *StateSet) bool {
	return d.dict.Contains(ss.bs)
}

func (d *StateSetDict[T]) Size() int {
	return len(d.shaToSs)
}

func (d *StateSetDict[T]) iterator() *stateSetDictIterator[T] {
	stSets := make([]*StateSet, 0)
	for _, ss := range d.shaToSs {
		stSets = append(stSets, ss)
	}

	return &stateSetDictIterator[T]{
		d:       d,
		stSets:  stSets,
		length:  len(stSets),
		currIdx: 0,
	}
}

type stateSetDictIterator[T any] struct {
	d       *StateSetDict[T]
	stSets  []*StateSet
	length  int
	currIdx int
}

func (iter *stateSetDictIterator[T]) HasNext() bool {
	return iter.currIdx < iter.length
}

func (iter *stateSetDictIterator[T]) Next() (*StateSet, T) {
	ss := iter.stSets[iter.currIdx]
	v, _ := iter.d.Get(ss)
	iter.currIdx++
	return ss, v
}

type stateUnionFind struct {
	uf *collection.UnionFind
}

func newStateUnionFind(n int) *stateUnionFind {
	return &stateUnionFind{
		uf: collection.NewUnionFind(n),
	}
}

func (uf *stateUnionFind) Unite(x, y StateID) bool {
	return uf.uf.Unite(int(x), int(y))
}

func (uf *stateUnionFind) Find(x StateID) StateID {
	return StateID(uf.uf.Find(int(x)))
}

func (uf *stateUnionFind) Same(x, y StateID) bool {
	return uf.Find(x) == uf.Find(y)
}

func charLabel(s string) string {
	switch s {
	case "\t":
		return "Tab"
	case "\n":
		return "Newline"
	case "\r":
		return "CR"
	case " ":
		return "Space"
	case "\\":
		return "Backslash"
	default:
		return s
	}
}
