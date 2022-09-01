package automata

import (
	"crypto/sha256"
	stdmath "math"

	"github.com/goropikari/golex/collection"
)

const epsilon = 'ε'

const SupportedChars = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ \t\n\r"

// const SupportedChars = "abcdefghijklmnopqrstuvwxyz+-*/.0123456789"

// const SupportedChars = "abc"

type TokenID int
type StateID int
type Sha = [sha256.Size]byte

type State struct {
	id      StateID
	tokenID TokenID
}

func NewState(id StateID) State {
	return State{id: id, tokenID: TokenID(stdmath.MaxInt)}
}

func (st State) GetID() StateID {
	return st.id
}

func (st State) GetTokenID() TokenID {
	if int(st.tokenID) == stdmath.MaxInt {
		return 0
	}
	return st.tokenID
}

func (st State) GetRawTokenID() TokenID {
	return st.tokenID
}

func (st *State) SetTokenID(id TokenID) {
	st.tokenID = id
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

func (ss *StateSet) IsEmpty() bool {
	return ss.bs.IsZero()
}

func (ss *StateSet) Contains(x StateID) bool {
	return ss.bs.Contains(int(x))
}

func (ss *StateSet) Sha256() Sha {
	return sha256.Sum256(ss.bs.Bytes())
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
	return newStateSetIterator(ss)
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
