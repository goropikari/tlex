package automata

import (
	"github.com/goropikari/golex/collection"
)

const blackHoleStateID = 0

type DFATransition struct {
	mp map[collection.Pair[State, rune]]State
}

func NewDFATransition() *DFATransition {
	return &DFATransition{
		mp: map[collection.Pair[State, rune]]State{},
	}
}

func (trans *DFATransition) Set(from State, ru rune, to State) {
	trans.mp[collection.NewPair(from, ru)] = to
}

func (trans *DFATransition) Erase(pair collection.Pair[State, rune]) {
	delete(trans.mp, pair)
}

func (trans *DFATransition) step(st State, ru rune) (State, bool) {
	st, ok := trans.mp[collection.NewPair(st, ru)]
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
	pairs := make([]collection.Pair[State, rune], 0, len(trans.mp))
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
	pairs   []collection.Pair[State, rune]
	tos     []State
}

func (iter *DFATransitionIterator) HasNext() bool {
	return iter.currIdx < iter.length
}

func (iter *DFATransitionIterator) Next() (collection.Pair[State, rune], State) {
	k := iter.pairs[iter.currIdx]
	v := iter.tos[iter.currIdx]
	iter.currIdx++

	return k, v
}

type DFA struct {
	q         *collection.Set[State]
	delta     *DFATransition
	initState State
	finStates *collection.Set[State]
}

func NewDFA(q *collection.Set[State], delta *DFATransition, initState State, finStates *collection.Set[State]) DFA {
	return DFA{
		q:         q,
		delta:     delta,
		initState: initState,
		finStates: finStates,
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

// func (dfa DFA) ToNFA() NFA {
// 	dfa = dfa.Copy().Minimize()
// 	delta := make(NFATransition)
// 	for pair, to := range dfa.delta {
// 		delta[pair] = collection.NewSet[State]().Insert(to)
// 	}

// 	return NewNFA(dfa.q, delta, collection.NewSet[State]().Insert(dfa.initState), dfa.finStates)
// }

func (dfa DFA) Accept(s string) (RegexID, bool) {
	currSt := dfa.initState

	for _, ru := range []rune(s) {
		var ok bool
		currSt, ok = dfa.Step(currSt, ru)
		if !ok { // implicit black hole state
			return 0, false
		}
		if currSt.GetID() == blackHoleStateID {
			return 0, false
		}
	}

	return currSt.GetRegexID(), dfa.finStates.Contains(currSt)
}

func (dfa DFA) Step(st State, ru rune) (State, bool) {
	return dfa.delta.step(st, ru)
}

func (dfa DFA) Copy() DFA {
	return NewDFA(dfa.q.Copy(), dfa.delta.Copy(), dfa.initState, dfa.finStates.Copy())
}

func (dfa DFA) Totalize() DFA {
	dfa = dfa.Copy()
	bhState := NewState(blackHoleStateID)
	states := dfa.q.Copy().Insert(bhState)
	delta := dfa.delta.Copy()
	changed := false
	for _, ru := range SupportedChars {
		qiter := dfa.q.Iterator()
		for qiter.HasNext() {
			st := qiter.Next()
			if _, ok := dfa.delta.step(st, ru); !ok {
				changed = true
				delta.Set(st, ru, bhState)
			}
		}
	}

	if changed {
		return NewDFA(states, delta, dfa.initState, dfa.finStates)
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
		ru := pair.Second
		tu := collection.NewPair(ns, ru)
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

const asciiSize = 255

func (dfa DFA) transTable() [][asciiSize]StateID {
	trans := make([][asciiSize]StateID, dfa.q.Size()+1)
	iter := dfa.q.Iterator()
	for iter.HasNext() {
		from := iter.Next()
		for _, ru := range SupportedChars {
			if to, ok := dfa.delta.step(from, ru); ok {
				trans[from.GetID()][ru] = to.GetID()
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
	numst := dfa.q.Size()
	qiter := dfa.q.Iterator()
	regSts := make(map[RegexID][]StateID)
	stIDToState := make(map[StateID]State)
	for qiter.HasNext() {
		st := qiter.Next()
		stID := st.GetID()
		regID := st.GetRawRegexID()
		regSts[regID] = append(regSts[regID], stID)
		stIDToState[stID] = st
	}

	grps := make([][]StateID, 0)
	uf := newStateGrouping(numst)
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
		newuf := newStateGrouping(numst)

		for _, group := range grps {
			for i, s0 := range group {
				for _, s1 := range group[i+1:] {
					same := true
					for _, ru := range SupportedChars {
						ns0 := transTable[s0][ru]
						ns1 := transTable[s1][ru]

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

func (dfa DFA) LexerMinimize() DFA {
	dfa = dfa.Totalize()
	groups := dfa.grouping()
	stIDToState := make(map[StateID]State)
	sts := dfa.q.Slice()
	for _, st := range sts {
		stIDToState[st.GetID()] = st
	}

	uf := newStateGrouping(dfa.q.Size())
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
	for _, st := range sts {
		q.Insert(stIDToState[uf.Find(st.GetID())])
	}

	initState := stIDToState[uf.Find(dfa.initState.GetID())]

	delta := NewDFATransition()
	iter := dfa.delta.Iterator()
	for iter.HasNext() {
		pair, ns := iter.Next()
		from := stIDToState[uf.Find(pair.First.GetID())]
		ru := pair.Second
		ns = stIDToState[uf.Find(ns.GetID())]
		delta.Set(from, ru, ns)
	}

	finStates := collection.NewSet[State]()
	fiter := dfa.finStates.Iterator()
	for fiter.HasNext() {
		st := fiter.Next()
		finStates.Insert(stIDToState[uf.Find(st.GetID())])
	}

	return NewDFA(q, delta, initState, finStates)
}
