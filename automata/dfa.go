package automata

import (
	"github.com/goropikari/tlex/collection"
)

const blackHoleStateID = 0

type DFATransition struct {
	mp map[collection.Pair[State, byte]]State
}

func NewDFATransition() *DFATransition {
	return &DFATransition{
		mp: map[collection.Pair[State, byte]]State{},
	}
}

func (trans *DFATransition) Set(from State, b byte, to State) {
	trans.mp[collection.NewPair(from, b)] = to
}

func (trans *DFATransition) Erase(pair collection.Pair[State, byte]) {
	delete(trans.mp, pair)
}

func (trans *DFATransition) step(st State, b byte) (State, bool) {
	st, ok := trans.mp[collection.NewPair(st, b)]
	return st, ok
}

func (t *DFATransition) Copy() *DFATransition {
	delta := NewDFATransition()
	iter := t.Iterator()
	for iter.HasNext() {
		p, to := iter.Next()
		delta.Set(p.First, p.Second, to)
	}

	return delta
}

func (trans *DFATransition) Iterator() *DFATransitionIterator {
	pairs := make([]collection.Pair[State, byte], 0, len(trans.mp))
	tos := make([]State, 0, len(pairs))
	for k, v := range trans.mp {
		pairs = append(pairs, k)
		tos = append(tos, v)
	}

	return &DFATransitionIterator{
		currIdx: 0,
		length:  len(pairs),
		pairs:   pairs,
		tos:     tos,
	}
}

type DFATransitionIterator struct {
	currIdx int
	length  int
	pairs   []collection.Pair[State, byte]
	tos     []State
}

func (iter *DFATransitionIterator) HasNext() bool {
	return iter.currIdx < iter.length
}

func (iter *DFATransitionIterator) Next() (collection.Pair[State, byte], State) {
	k := iter.pairs[iter.currIdx]
	v := iter.tos[iter.currIdx]
	iter.currIdx++

	return k, v
}

type DFA struct {
	q             *collection.Set[State]
	delta         *DFATransition
	initState     State
	finStates     *collection.Set[State]
	stIDToRegexID StateIDToRegexID
}

func NewDFA(q *collection.Set[State], delta *DFATransition, initState State, finStates *collection.Set[State], stIDToRegexID StateIDToRegexID) DFA {
	return DFA{
		q:             q,
		delta:         delta,
		initState:     initState,
		finStates:     finStates,
		stIDToRegexID: stIDToRegexID,
	}
}

func (dfa DFA) GetStates() []State {
	return dfa.q.Slice()
}

func (dfa DFA) GetInitState() State {
	return dfa.initState
}

func (dfa DFA) GetFinStates() *collection.Set[State] {
	return dfa.finStates
}

func (dfa DFA) GetTransitionTable() *DFATransition {
	return dfa.delta
}

func (dfa DFA) GetRegexID(st State) RegexID {
	return dfa.stIDToRegexID.Get(st.GetID())
}

func (dfa DFA) Accept(s string) (RegexID, bool) {
	currSt := dfa.initState

	for _, b := range []byte(s) {
		var ok bool
		currSt, ok = dfa.Step(currSt, b)
		if !ok { // implicit black hole state
			return 0, false
		}
		if currSt.GetID() == blackHoleStateID {
			return 0, false
		}
	}

	return dfa.stIDToRegexID.Get(currSt.GetID()), dfa.finStates.Contains(currSt)
}

func (dfa DFA) Step(st State, b byte) (State, bool) {
	return dfa.delta.step(st, b)
}

func (dfa DFA) Copy() DFA {
	return NewDFA(dfa.q.Copy(), dfa.delta.Copy(), dfa.initState, dfa.finStates.Copy(), dfa.stIDToRegexID)
}

func (dfa DFA) Totalize() DFA {
	dfa = dfa.Copy()
	bhState := NewState(blackHoleStateID)
	states := dfa.q.Copy().Insert(bhState)
	delta := dfa.delta.Copy()
	changed := false
	for _, b := range SupportedChars {
		qiter := dfa.q.Iterator()
		for qiter.HasNext() {
			st := qiter.Next()
			if _, ok := dfa.delta.step(st, b); !ok {
				changed = true
				delta.Set(st, b, bhState)
			}
		}
	}

	if changed {
		return NewDFA(states, delta, dfa.initState, dfa.finStates, dfa.stIDToRegexID)
	}

	return dfa
}

func (dfa DFA) Reverse() NFA {
	dfa = dfa.Totalize()
	delta := make(NFATransition)
	iter := dfa.delta.Iterator()
	for iter.HasNext() {
		pair, ns := iter.Next()
		from := pair.First
		b := pair.Second
		tu := collection.NewPair(ns, b)
		if _, ok := delta[tu]; ok {
			delta[tu].Insert(from)
		} else {
			delta[tu] = collection.NewSet[State]().Insert(from)
		}
	}

	return NewNFA(dfa.q, delta, dfa.finStates, collection.NewSet[State]().Insert(dfa.initState))
}

// // Brzozowski DFA minimization algorithm
// func (dfa DFA) Minimize() DFA {
// 	return dfa.Reverse().ToDFA().Reverse().ToDFA()
// }

func (dfa DFA) RemoveBH() DFA {
	dfa = dfa.Copy()

	bhSt := NewState(blackHoleStateID)
	dfa.q.Erase(bhSt)

	iter := dfa.delta.Iterator()
	for iter.HasNext() {
		pair, to := iter.Next()
		if to.GetID() == blackHoleStateID {
			dfa.delta.Erase(pair)
		}
	}

	return dfa
}

type stateGroup struct {
	states *collection.Set[State]
}

func NewGroup(states *collection.Set[State]) *stateGroup {
	return &stateGroup{states: states}
}

func (g *stateGroup) size() int {
	return g.states.Size()
}

func (g *stateGroup) slice() []State {
	sts := make([]State, 0)
	iter := g.states.Iterator()
	for iter.HasNext() {
		st := iter.Next()
		sts = append(sts, st)
	}

	return sts
}

func (dfa DFA) transTable() [][asciiSize]StateID {
	trans := make([][asciiSize]StateID, dfa.q.Size()+1)
	iter := dfa.q.Iterator()
	for iter.HasNext() {
		from := iter.Next()
		for _, b := range SupportedChars {
			if to, ok := dfa.delta.step(from, b); ok {
				trans[from.GetID()][b] = to.GetID()
			}
		}
	}

	return trans
}

// state minimization for lexical analyzer
// Compilers: Principles, Techniques, and Tools, 2ed ed.,  ISBN 9780321486813 (Dragon book)
// p.181 Algorithm 3.39
// p.184 3.9.7 State Minimization in Lexical Analyzers
func (dfa DFA) grouping() []*stateGroup {
	numst := dfa.q.Size() + 1
	qiter := dfa.q.Iterator()
	regSts := make(map[RegexID][]StateID)
	stIDToState := make(map[StateID]State)
	for qiter.HasNext() {
		st := qiter.Next()
		stID := st.GetID()
		regID := dfa.stIDToRegexID.Get(stID)
		regSts[regID] = append(regSts[regID], stID)
		stIDToState[stID] = st
	}

	grps := make([][]StateID, 0)
	uf := newStateUnionFind(numst)
	for _, stIDs := range regSts {
		for _, stID := range stIDs[1:] {
			uf.Unite(stIDs[0], stID)
		}
		grps = append(grps, stIDs)
	}

	transTable := dfa.transTable()

	ngrp := len(grps)
	splitted := true
	for splitted {
		splitted = false
		newuf := newStateUnionFind(numst)

		for _, group := range grps {
			for i, s0 := range group {
				for _, s1 := range group[i+1:] {
					same := true
					for _, b := range SupportedChars {
						ns0 := transTable[s0][b]
						ns1 := transTable[s1][b]

						if !uf.Same(ns0, ns1) {
							same = false
							break
						}
					}
					if same {
						newuf.Unite(s0, s1)
					}
				}
			}
		}

		mp := make(map[StateID][]StateID)
		for stID := StateID(0); stID < StateID(numst); stID++ {
			mp[newuf.Find(stID)] = append(mp[newuf.Find(stID)], stID)
		}
		newGrps := make([][]StateID, 0, len(mp))
		for _, v := range mp {
			newGrps = append(newGrps, v)
		}

		uf = newuf
		splitted = ngrp != len(newGrps)
		ngrp = len(newGrps)
		grps = newGrps
	}

	stGrps := make([]*stateGroup, 0, len(grps))
	for _, grp := range grps {
		sets := collection.NewSet[State]()
		for _, stID := range grp {
			sets.Insert(stIDToState[stID])
		}
		stGrps = append(stGrps, NewGroup(sets))
	}

	return stGrps
}

// minimize したあとの StateID は連番にはなっていない
func (dfa DFA) LexerMinimize() DFA {
	dfa = dfa.Totalize()
	groups := dfa.grouping()
	stIDToState := make(map[StateID]State)
	sts := dfa.q.Slice()
	for _, st := range sts {
		stIDToState[st.GetID()] = st
	}

	uf := newStateUnionFind(dfa.q.Size() + 1)
	for _, g := range groups {
		n := g.size()
		if n == 1 {
			continue
		}
		sts := g.slice()
		for i := 1; i < n; i++ {
			uf.Unite(sts[0].GetID(), sts[i].GetID())
		}
	}

	q := collection.NewSet[State]()
	stIDToRegexID := make(StateIDToRegexID)
	for _, st := range sts {
		leaderID := uf.Find(st.GetID())
		q.Insert(stIDToState[leaderID])
		stIDToRegexID.Set(leaderID, dfa.stIDToRegexID.Get(leaderID))
	}

	initState := stIDToState[uf.Find(dfa.initState.GetID())]

	delta := NewDFATransition()
	iter := dfa.delta.Iterator()
	for iter.HasNext() {
		pair, ns := iter.Next()
		from := stIDToState[uf.Find(pair.First.GetID())]
		b := pair.Second
		ns = stIDToState[uf.Find(ns.GetID())]
		delta.Set(from, b, ns)
	}

	finStates := collection.NewSet[State]()
	fiter := dfa.finStates.Iterator()
	for fiter.HasNext() {
		st := fiter.Next()
		finStates.Insert(stIDToState[uf.Find(st.GetID())])
	}

	return NewDFA(q, delta, initState, finStates, stIDToRegexID)
}
