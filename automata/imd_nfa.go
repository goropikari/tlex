package automata

import (
	"errors"

	"github.com/goropikari/tlex/collection"
	"github.com/goropikari/tlex/math"
)

type ImdEpsilonTransition struct {
	size int
	mp   map[StateID]*StateSet
}

func NewImdEpsilonTransition(size int, mp map[StateID]*StateSet) ImdEpsilonTransition {
	return ImdEpsilonTransition{
		size: size,
		mp:   mp,
	}
}

func (trans *ImdEpsilonTransition) step(from StateID) *StateSet {
	return trans.mp[from]
}

type ImdNFATransition struct {
	mp map[StateID]map[Interval]*StateSet
}

func NewImdNFATransition(mp map[StateID]map[Interval]*StateSet) ImdNFATransition {
	return ImdNFATransition{mp: mp}
}

type ImdNFA struct {
	size        int
	intvs       []Interval
	etrans      ImdEpsilonTransition
	trans       ImdNFATransition
	initStates  *StateSet
	finStates   *StateSet
	stIDToRegID StateIDToRegexID
}

func (nfa *ImdNFA) ToDFA() *DFA {
	stateSetDict, imdTrans, imdInitState, imdFinStates := nfa.SubsetConstruction()

	states := collection.NewSet[StateID]()
	initState, ok := stateSetDict.Get(imdInitState)
	if !ok {
		panic(errors.New("cannot find initial state"))
	}
	states.Insert(initState)

	trans := NewDFATransition()
	titer := imdTrans.iterator()
	for titer.HasNext() {
		fromSs, mp := titer.Next()
		fsid, ok := stateSetDict.Get(fromSs)
		if !ok {
			panic(errors.New("cannot find given state"))
		}
		for intv, toss := range mp {
			toid, ok := stateSetDict.Get(toss)
			if !ok {
				panic(errors.New("cannot find given state"))
			}
			trans.Set(fsid, intv, toid)
			states.Insert(toid)
		}
	}

	stIDToRegID := NewStateIDToRegexID()
	finStates := collection.NewSet[StateID]()
	fiter := imdFinStates.iterator()
	for fiter.HasNext() {
		ss, _ := fiter.Next()
		sid, ok := stateSetDict.Get(ss)
		if !ok {
			panic(errors.New("cannot find given state"))
		}
		finStates.Insert(sid)

		regID := nonFinStateRegexID
		siter := ss.iterator()
		for siter.HasNext() {
			rid := nfa.stIDToRegID.Get(siter.Next())
			regID = math.Min(regID, rid)
		}
		stIDToRegID.Set(sid, regID)
	}

	return &DFA{
		size:        stateSetDict.Size(),
		states:      states,
		intvs:       nfa.intvs,
		trans:       trans,
		initState:   initState,
		finStates:   finStates,
		stIDToRegID: stIDToRegID,
	}
}

func (nfa *ImdNFA) SubsetConstruction() (states *StateSetDict[StateID], trans *ImdDFATransition, initState *StateSet, finStates *StateSetDict[Nothing]) {
	n := nfa.size
	ecls := make([]*StateSet, n)
	for sid := StateID(0); sid < StateID(n); sid++ {
		ecls[sid] = nfa.Eclosure(sid)
	}

	initState = NewStateSet(n)
	iiter := nfa.initStates.iterator()
	for iiter.HasNext() {
		sid := iiter.Next()
		initState = initState.Union(ecls[sid])
	}

	visited := NewStateSetDict[StateID]() // key is set of states, value is state id for the key.
	finStateDict := NewStateSetDict[Nothing]()
	if initState.Intersection(nfa.finStates).IsAny() {
		finStateDict.Set(initState, nothing)
	}

	delta := NewImdDFATransition()
	deq := collection.NewDeque[*StateSet]()
	deq.PushBack(initState)

	id := StateID(0)
	for deq.Size() > 0 {
		froms := deq.Front()
		deq.PopFront()

		if visited.Contains(froms) {
			continue
		}
		visited.Set(froms, id)
		id++

		// nfa.intvs is already disjoined.
		for _, intv := range nfa.intvs {
			tos := NewStateSet(n)
			fiter := froms.iterator()
			for fiter.HasNext() {
				fsid := fiter.Next()
				nxs := nfa.step(fsid, intv)
				niter := nxs.iterator()
				for niter.HasNext() {
					nsid := niter.Next()
					tos = tos.Union(ecls[nsid])
				}
			}

			if tos.IsEmpty() {
				continue
			}

			if tos.Intersection(nfa.finStates).IsAny() {
				finStateDict.Set(tos, nothing)
			}

			delta.Set(froms, intv, tos)

			if visited.Contains(tos) {
				continue
			}
			deq.PushBack(tos)
		}
	}

	return visited, delta, initState, finStateDict
}

func (nfa *ImdNFA) step(fsid StateID, intv Interval) *StateSet {
	ss := NewStateSet(nfa.size)

	for fintv, nxs := range nfa.trans.mp[fsid] {
		if fintv.Overlap(intv) {
			ss = ss.Union(nxs)
		}
	}

	return ss
}

func (nfa *ImdNFA) Eclosure(sid StateID) *StateSet {
	ss := NewStateSet(nfa.size)

	deq := collection.NewDeque[StateID]()
	deq.PushBack(sid)
	for deq.Size() > 0 {
		fr := deq.Front()
		deq.PopFront()

		if ss.Contains(fr) {
			continue
		}
		ss.Insert(fr)

		nxs := nfa.etrans.step(fr)
		iter := nxs.iterator()
		for iter.HasNext() {
			nx := iter.Next()
			if ss.Contains(nx) {
				continue
			}
			deq.PushBack(nx)
		}
	}

	return ss
}

type ImdDFATransition struct {
	d *StateSetDict[map[Interval]*StateSet]
}

func NewImdDFATransition() *ImdDFATransition {
	return &ImdDFATransition{
		d: NewStateSetDict[map[Interval]*StateSet](),
	}
}

func (trans *ImdDFATransition) Set(from *StateSet, intv Interval, to *StateSet) {
	if v, ok := trans.d.Get(from); ok {
		v[intv] = to
		trans.d.Set(from, v)
		return
	}

	mp := make(map[Interval]*StateSet)
	mp[intv] = to
	trans.d.Set(from, mp)
}

func (trans *ImdDFATransition) iterator() *stateSetDictIterator[map[Interval]*StateSet] {
	return trans.d.iterator()
}
