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

// state minimization for lexical analyzer
// Compilers: Principles, Techniques, and Tools, 2ed ed.,  ISBN 9780321486813 (Dragon book)
// p.181 Algorithm 3.39
// p.184 3.9.7 State Minimization in Lexical Analyzers
func (dfa DFA) grouping() []*stateGroup {
	states := dfa.q.Slice()

	stateSets := make(map[RegexID]*collection.Set[State])
	qiter := dfa.q.Iterator()
	for qiter.HasNext() {
		st := qiter.Next()
		if _, ok := stateSets[st.GetRegexID()]; ok {
			stateSets[st.GetRegexID()].Insert(st)
		} else {
			stateSets[st.GetRegexID()] = collection.NewSet[State]().Insert(st)
		}
	}
	groups := make([]*stateGroup, 0, len(stateSets))
	for _, group := range stateSets {
		groups = append(groups, NewGroup(group))
	}

	ngrp := len(groups)
	isSplit := true
	for isSplit {
		isSplit = false

		// old groups
		oldStUF := newStateUnionFind(states)
		for _, grp := range groups {
			gss := grp.slice()
			if len(gss) == 1 {
				continue
			}
			for i := 0; i < len(gss); i++ {
				oldStUF.unite(gss[0], gss[i])
			}
		}

		// new groups
		newStUF := newStateUnionFind(states)
		for _, grp := range groups {
			gss := grp.slice()
			if len(gss) == 1 {
				continue
			}
			for i := 0; i < len(gss); i++ {
				for j := i + 1; j < len(gss); j++ {
					s0 := gss[i]
					s1 := gss[j]
					isSameGroup := true
					for _, ru := range SupportedChars {
						ns0, _ := dfa.delta.step(s0, ru)
						ns1, _ := dfa.delta.step(s1, ru)
						// If ns0 and ns1 belong to different groups, s0 and s1 belong to other groups.
						// Then current group is split.
						if !oldStUF.same(ns0, ns1) {
							isSameGroup = false
							break
						}
					}
					if isSameGroup {
						newStUF.unite(s0, s1)
					}
				}
			}
		}

		newStateSets := make(map[State]*collection.Set[State])
		for _, st := range states {
			leaderSt := newStUF.find(st)
			if _, ok := newStateSets[leaderSt]; ok {
				newStateSets[leaderSt].Insert(st)
			} else {
				newStateSets[leaderSt] = collection.NewSet[State]().Insert(st)
			}
		}
		newGroups := make([]*stateGroup, 0)
		for _, group := range newStateSets {
			newGroups = append(newGroups, NewGroup(group))
		}

		// If group splitting occurs, the number of groups is increasing.
		if ngrp != len(newGroups) {
			ngrp = len(newGroups)
			isSplit = true
			groups = newGroups
		}
	}

	return groups
}

func (dfa DFA) LexerMinimize() DFA {
	dfa = dfa.Totalize()
	groups := dfa.grouping()
	states := dfa.q.Slice()

	uf := newStateUnionFind(states)
	for _, g := range groups {
		n := g.size()
		if n == 1 {
			continue
		}
		states := g.slice()
		for i := 1; i < n; i++ {
			uf.unite(states[0], states[i])
		}
	}

	q := collection.NewSet[State]()
	qiter := dfa.q.Iterator()
	for qiter.HasNext() {
		st := qiter.Next()
		q.Insert(uf.find(st))
	}

	initState := uf.find(dfa.initState)

	delta := NewDFATransition()
	iter := dfa.delta.Iterator()
	for iter.HasNext() {
		pair, ns := iter.Next()
		from := uf.find(pair.First)
		ru := pair.Second
		ns = uf.find(ns)
		delta.Set(from, ru, ns)
	}

	finStates := collection.NewSet[State]()
	fiter := dfa.finStates.Iterator()
	for fiter.HasNext() {
		st := fiter.Next()
		finStates.Insert(uf.find(st))
	}

	return NewDFA(q, delta, initState, finStates)
}
