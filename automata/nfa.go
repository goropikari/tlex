package automata

import (
	"github.com/goropikari/tlex/collection"
)

type EpsilonTransition struct {
	mp map[StateID]*collection.Set[StateID]
}

func NewEpsilonTransition(mp map[StateID]*collection.Set[StateID]) EpsilonTransition {
	return EpsilonTransition{
		mp: mp,
	}
}

func (t EpsilonTransition) set(from, to StateID) {
	if _, ok := t.mp[from]; ok {
		t.mp[from].Insert(to)
	} else {
		t.mp[from] = collection.NewSet[StateID]().Insert(to)
	}
}

func (trans *EpsilonTransition) merge(other EpsilonTransition) {
	for sid, set := range other.mp {
		if v, ok := trans.mp[sid]; ok {
			trans.mp[sid] = v.Union(set)
		} else {
			trans.mp[sid] = set
		}
	}
}

func (trans *EpsilonTransition) step(sid StateID) *collection.Set[StateID] {
	return trans.mp[sid]
}

type NFATransition struct {
	mp map[StateID]map[Interval]*collection.Set[StateID]
}

func NewTransition(mp map[StateID]map[Interval]*collection.Set[StateID]) NFATransition {
	return NFATransition{mp: mp}
}

func (trans NFATransition) merge(other NFATransition) {
	for sid, mp := range other.mp {
		trans.mp[sid] = mp
	}
}

func (trans NFATransition) step(sid StateID, intv Interval) *collection.Set[StateID] {
	mp, ok := trans.mp[sid]
	if !ok {
		return nil
	}

	return mp[intv]
}

func (trans NFATransition) intervals() []Interval {
	intvs := make([]Interval, 0)
	for _, mp := range trans.mp {
		for intv := range mp {
			intvs = append(intvs, intv)
		}
	}

	return Disjoin(intvs)
}

type NFA struct {
	states       *collection.Set[StateID]
	epsilonTrans EpsilonTransition
	trans        NFATransition
	initStates   *collection.Set[StateID]
	finStates    *collection.Set[StateID]
	stIDToRegID  StateIDToRegexID
}

func NewNFA(states *collection.Set[StateID], etrans EpsilonTransition, trans NFATransition, initStates *collection.Set[StateID], finStates *collection.Set[StateID]) *NFA {
	return &NFA{
		states:       states,
		epsilonTrans: etrans,
		trans:        trans,
		initStates:   initStates,
		finStates:    finStates,
	}
}

func (nfa *NFA) Sum(other *NFA) *NFA {
	nfa.states = nfa.states.Union(other.states)
	nfa.epsilonTrans.merge(other.epsilonTrans)
	nfa.trans.merge(other.trans)

	oiter := other.initStates.Iterator()
	for oiter.HasNext() {
		sid := oiter.Next()
		nfa.initStates.Insert(sid)
	}
	fiter := other.finStates.Iterator()
	for fiter.HasNext() {
		sid := fiter.Next()
		nfa.finStates.Insert(sid)
	}

	for sid, rid := range other.stIDToRegID {
		nfa.stIDToRegID.Set(sid, rid)
	}

	return nfa
}

func (nfa *NFA) Concat(other *NFA) *NFA {
	nfa.states = nfa.states.Union(other.states)
	nfa.epsilonTrans.merge(other.epsilonTrans)
	nfa.trans.merge(other.trans)
	fiter := nfa.finStates.Iterator()
	for fiter.HasNext() {
		from := fiter.Next()
		iiter := other.initStates.Iterator()
		for iiter.HasNext() {
			to := iiter.Next()
			nfa.epsilonTrans.set(from, to)
		}
	}
	nfa.finStates = other.finStates

	return nfa
}

func (nfa *NFA) Star() *NFA {
	sid := NewStateID()

	nfa.states = nfa.states.Insert(sid)
	fiter := nfa.finStates.Iterator()
	for fiter.HasNext() {
		from := fiter.Next()
		nfa.epsilonTrans.set(from, sid)
	}
	iiter := nfa.initStates.Iterator()
	for iiter.HasNext() {
		to := iiter.Next()
		nfa.epsilonTrans.set(sid, to)
	}

	states := collection.NewSet[StateID]().Insert(sid)
	nfa.initStates = states
	nfa.finStates = states

	return nfa
}

func (nfa *NFA) SetRegexID(rid RegexID) *NFA {
	if nfa.stIDToRegID == nil {
		nfa.stIDToRegID = make(StateIDToRegexID)
	}
	iter := nfa.finStates.Iterator()
	for iter.HasNext() {
		sid := iter.Next()
		nfa.stIDToRegID.Set(sid, rid)
	}

	return nfa
}

func (nfa *NFA) ToImdNFA() *ImdNFA {
	n := nfa.states.Size()
	oldToNew := map[StateID]StateID{}
	newsid := StateID(0)
	siter := nfa.states.Iterator()
	for siter.HasNext() {
		oldsid := siter.Next()
		oldToNew[oldsid] = newsid
		newsid++
	}

	epsilonMap := map[StateID]*StateSet{}
	for oldsid, newsid := range oldToNew {
		tos := nfa.epsilonTrans.step(oldsid)
		ss := NewStateSet(n)
		iter := tos.Iterator()
		for iter.HasNext() {
			nsid := iter.Next()
			ss.Insert(oldToNew[nsid])
		}
		epsilonMap[newsid] = ss
	}

	trans := map[StateID]map[Interval]*StateSet{}
	for sid, mp := range nfa.trans.mp {
		trans[oldToNew[sid]] = map[Interval]*StateSet{}
		for intv, nxs := range mp {
			ss := NewStateSet(n)
			iter := nxs.Iterator()
			for iter.HasNext() {
				nsid := iter.Next()
				ss.Insert(oldToNew[nsid])
			}
			trans[oldToNew[sid]][intv] = ss
		}
	}

	initStates := NewStateSet(n)
	iiter := nfa.initStates.Iterator()
	for iiter.HasNext() {
		sid := iiter.Next()
		initStates.Insert(oldToNew[sid])
	}

	finStates := NewStateSet(n)
	fiter := nfa.finStates.Iterator()
	for fiter.HasNext() {
		sid := fiter.Next()
		finStates.Insert(oldToNew[sid])
	}

	stIDToRegID := StateIDToRegexID{}
	for sid, rid := range nfa.stIDToRegID {
		stIDToRegID.Set(oldToNew[sid], rid)
	}

	return &ImdNFA{
		size:        n,
		intvs:       nfa.trans.intervals(),
		etrans:      NewImdEpsilonTransition(n, epsilonMap),
		trans:       NewImdNFATransition(trans),
		initStates:  initStates,
		finStates:   finStates,
		stIDToRegID: stIDToRegID,
	}
}
