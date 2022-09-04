package automata

import (
	"container/list"

	"github.com/goropikari/tlex/collection"
	"github.com/goropikari/tlex/math"
	"github.com/goropikari/tlex/utils/counter"
)

type ImdNFATransition map[collection.Pair[StateID, byte]]*StateSet

func (trans ImdNFATransition) step(x StateID, b byte) (*StateSet, bool) {
	nxs, ok := trans[collection.NewPair(x, b)]
	return nxs, ok
}

type ImdNFA struct {
	cnt         *counter.Counter
	maxID       int
	stIDToRegID []RegexID
	delta       ImdNFATransition
	initStates  *StateSet
	finStates   *StateSet
}

func NewImdNFA(maxID int, stIDToRegID []RegexID, delta ImdNFATransition, initStates *StateSet, finStates *StateSet) ImdNFA {
	return ImdNFA{
		cnt:         counter.NewCounter(1),
		maxID:       maxID,
		stIDToRegID: stIDToRegID,
		delta:       delta,
		initStates:  initStates,
		finStates:   finStates,
	}
}

func (nfa ImdNFA) buildEClosures() []*StateSet {
	ecl := make([]*StateSet, nfa.numst())
	iter := nfa.iterator()
	for iter.HasNext() {
		sid := iter.Next()
		b := nfa.eclosure(sid)
		ecl[sid] = b
	}

	return ecl
}

func (nfa ImdNFA) numst() int {
	// +1 means black hole
	return nfa.maxID + 1
}

func (nfa ImdNFA) genStateID() StateID {
	return StateID(nfa.cnt.Generate())
}

func (nfa ImdNFA) calRegID(ss *StateSet) RegexID {
	regID := nonFinStateRegexID
	iter := ss.iterator()
	for iter.HasNext() {
		sid := iter.Next()
		regID = math.Min(regID, nfa.stIDToRegID[sid])
	}

	return regID
}

func (nfa ImdNFA) step(sid StateID, b byte) (*StateSet, bool) {
	nxid, ok := nfa.delta.step(sid, b)
	return nxid, ok
}

func (nfa ImdNFA) ToDFA() DFA {
	states, delta, initState, finStates := nfa.subsetConstruction()

	dfaStates := collection.NewSet[State]()
	dfaStIDToRegexID := make(map[StateID]RegexID)
	ssToSt := NewStateSetDict[State]()
	siter := states.iterator()
	for siter.HasNext() {
		ss, newSid := siter.Next()
		regID := nfa.calRegID(ss)
		dfaStIDToRegexID[newSid] = regID
		st := NewState(newSid)
		dfaStates.Insert(st)
		ssToSt.Set(ss, st)
	}

	diter := delta.iterator()
	dfaDelta := NewDFATransition()
	for diter.HasNext() {
		fromSs, mp := diter.Next()
		for b, toSs := range mp {
			fromSt, _ := ssToSt.Get(fromSs)
			toSt, _ := ssToSt.Get(toSs)
			dfaDelta.Set(fromSt, b, toSt)
		}
	}

	dfaInitState, _ := ssToSt.Get(initState)

	dfaFinStates := collection.NewSet[State]()
	fiter := finStates.iterator()
	for fiter.HasNext() {
		ss, _ := fiter.Next()
		st, _ := ssToSt.Get(ss)
		dfaFinStates.Insert(st)
	}

	return NewDFA(dfaStates, dfaDelta, dfaInitState, dfaFinStates, dfaStIDToRegexID)
}

func (nfa ImdNFA) subsetConstruction() (states *StateSetDict[StateID], delta *StateSetDict[map[byte]*StateSet], initState *StateSet, finStates *StateSetDict[Nothing]) {
	ecl := nfa.buildEClosures()

	initState = nfa.initStates.Copy()
	initIter := initState.iterator()
	for initIter.HasNext() {
		sid := initIter.Next()
		initState = initState.Union(ecl[sid])
	}

	visited := NewStateSetDict[StateID]()
	finStates = NewStateSetDict[Nothing]()
	if initState.Intersection(nfa.finStates).IsAny() {
		finStates.Set(initState, nothing)
	}
	visited.Set(initState, nfa.genStateID())

	delta = NewStateSetDict[map[byte]*StateSet]()

	que := list.New() // list of *StateSet
	que.PushBack(initState)
	for que.Len() > 0 {
		top := que.Front()
		que.Remove(top)
		from := top.Value.(*StateSet)

		for _, b := range SupportedChars {
			tos := NewStateSet(nfa.numst())
			fromIter := from.iterator()
			for fromIter.HasNext() {
				fromStID := fromIter.Next()
				if nxs, ok := nfa.step(fromStID, b); ok {
					nxsIter := nxs.iterator()
					for nxsIter.HasNext() {
						nxStID := nxsIter.Next()
						tos = tos.Union(ecl[nxStID])
					}
				}
			}

			if tos.IsEmpty() {
				continue
			}
			if tos.Intersection(nfa.finStates).IsAny() {
				finStates.Set(tos, nothing)
			}
			if v, ok := delta.Get(from); ok {
				v[b] = tos
				delta.Set(from, v)
			} else {
				mp := map[byte]*StateSet{}
				mp[b] = tos
				delta.Set(from, mp)
			}
			if visited.Contains(tos) {
				continue
			}
			visited.Set(tos, nfa.genStateID())
			que.PushBack(tos)
		}
	}

	return visited, delta, initState, finStates
}

func (nfa ImdNFA) eclosure(x StateID) *StateSet {
	que := list.New() // list of StateID
	que.PushBack(x)

	visited := NewStateSet(nfa.maxID + 1).Insert((x))
	closure := visited.Copy()
	for que.Len() > 0 {
		front := que.Front()
		que.Remove(front)
		top := front.Value.(StateID)

		if nxs, ok := nfa.step(top, epsilon); ok {
			closure = closure.Union(nxs)
			nxsIter := nxs.iterator()
			for nxsIter.HasNext() {
				nxStID := nxsIter.Next()
				if !visited.Contains(nxStID) {
					visited = visited.Insert(nxStID)
					que.PushBack(nxStID)
				}
			}
		}
	}

	return closure
}

func (nfa ImdNFA) iterator() *allStateIDIterator {
	return newAllStateIDIterator(nfa.maxID)
}

type allStateIDIterator struct {
	maxID  int
	currID int
}

func newAllStateIDIterator(maxID int) *allStateIDIterator {
	return &allStateIDIterator{
		maxID:  maxID,
		currID: 1, // StateID = 0 is blackhole state
	}
}

func (iter *allStateIDIterator) HasNext() bool {
	return iter.currID <= iter.maxID
}

func (iter *allStateIDIterator) Next() StateID {
	ret := StateID(iter.currID)
	iter.currID++
	return ret
}
