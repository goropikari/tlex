package automata

import (
	"github.com/goropikari/tlex/collection"
	"github.com/goropikari/tlex/utils/guid"
)

type NFATransition map[collection.Pair[State, byte]]*collection.Set[State]

func (t NFATransition) Copy() NFATransition {
	delta := make(NFATransition)
	for k, v := range t {
		delta[k] = v.Copy()
	}

	return delta
}

type NFA struct {
	q             *collection.Set[State]
	delta         NFATransition
	initStates    *collection.Set[State]
	finStates     *collection.Set[State]
	stIDToRegexID StateIDToRegexID
}

func NewNFA(
	q *collection.Set[State],
	delta NFATransition,
	initStates *collection.Set[State],
	finStates *collection.Set[State]) NFA {
	return NFA{
		q: q,
		// sigma:      sigma,
		delta:         delta,
		initStates:    initStates,
		finStates:     finStates,
		stIDToRegexID: make(StateIDToRegexID),
	}
}

func NewNFAWithRegexIDMap(q *collection.Set[State], delta NFATransition, initState *collection.Set[State], finState *collection.Set[State], stIDToRegexID map[StateID]RegexID) NFA {
	return NFA{
		q:             q,
		delta:         delta,
		initStates:    initState,
		finStates:     finState,
		stIDToRegexID: stIDToRegexID,
	}
}

func (nfa NFA) Copy() NFA {
	return NewNFAWithRegexIDMap(nfa.q.Copy(), nfa.delta.Copy(), nfa.initStates.Copy(), nfa.finStates.Copy(), nfa.stIDToRegexID)
}

func (nfa NFA) Concat(other NFA) NFA {
	nfa = nfa.Copy()
	other = other.Copy()

	qiter := other.q.Iterator()
	for qiter.HasNext() {
		st := qiter.Next()
		nfa.q.Insert(st)
	}

	for tr, ss := range other.delta {
		nfa.delta[tr] = ss
	}

	fiter := nfa.finStates.Iterator()
	for fiter.HasNext() {
		from := fiter.Next()
		iiter := other.initStates.Iterator()
		for iiter.HasNext() {
			to := iiter.Next()
			if _, ok := nfa.delta[collection.NewPair(from, epsilon)]; ok {
				nfa.delta[collection.NewPair(from, epsilon)].Insert(to)
			} else {
				nfa.delta[collection.NewPair(from, epsilon)] = collection.NewSet[State]().Insert(to)
			}
		}

	}

	return NewNFA(nfa.q, nfa.delta, nfa.initStates, other.finStates)
}

func (nfa NFA) Sum(other NFA) NFA {
	nfa = nfa.Copy()
	other = other.Copy()

	qiter := other.q.Iterator()
	for qiter.HasNext() {
		st := qiter.Next()
		nfa.q.Insert(st)
	}

	for tr, ss := range other.delta {
		nfa.delta[tr] = ss
	}

	iiter := other.initStates.Iterator()
	for iiter.HasNext() {
		st := iiter.Next()
		nfa.initStates.Insert(st)
	}

	fiter := other.finStates.Iterator()
	for fiter.HasNext() {
		st := fiter.Next()
		nfa.finStates.Insert(st)
	}

	return NewNFA(nfa.q, nfa.delta, nfa.initStates, nfa.finStates)
}

func (nfa NFA) SumWithRegexID(other NFA) NFA {
	stIDToRegexID := make(StateIDToRegexID)
	for k, v := range nfa.stIDToRegexID {
		stIDToRegexID.Set(k, v)
	}
	for k, v := range other.stIDToRegexID {
		stIDToRegexID.Set(k, v)
	}

	nfa = nfa.Sum(other)

	return NewNFAWithRegexIDMap(nfa.q, nfa.delta, nfa.initStates, nfa.finStates, stIDToRegexID)
}

func (nfa NFA) Star() NFA {
	nfa = nfa.Copy()

	startFinState := NewState(StateID(guid.New()))
	initStates := collection.NewSet[State]().Insert(startFinState)

	nfa.q.Insert(startFinState)

	nfa.delta[collection.NewPair(startFinState, epsilon)] = nfa.initStates

	fiter := nfa.finStates.Iterator()
	for fiter.HasNext() {
		from := fiter.Next()
		pair := collection.NewPair(from, epsilon)
		if _, ok := nfa.delta[pair]; ok {
			nfa.delta[pair].Insert(startFinState)
		} else {
			nfa.delta[pair] = initStates
		}
	}

	return NewNFA(nfa.q, nfa.delta, initStates, initStates)
}

func (nfa NFA) ToImNFA() ImdNFA {
	nfa = nfa.relabelStateIDs()
	maxID := nfa.q.Size()
	numst := maxID + 1 // +1 means black hole state
	stIDToRegID := make([]RegexID, numst)
	qiter := nfa.q.Iterator()
	for qiter.HasNext() {
		st := qiter.Next()
		sid := st.GetID()
		stIDToRegID[sid] = nfa.stIDToRegexID.Get(sid)
	}
	delta := make(ImdNFATransition)
	for pair, tos := range nfa.delta {
		from := pair.First
		b := pair.Second
		delta[collection.NewPair(from.GetID(), b)] = buildStateSet(numst, tos)
	}
	initStates := buildStateSet(numst, nfa.initStates)
	finStates := buildStateSet(numst, nfa.finStates)

	return NewImdNFA(maxID, stIDToRegID, delta, initStates, finStates)
}

func (nfa NFA) relabelStateIDs() NFA {
	nfa = nfa.Copy()

	id := StateID(1)
	oldToNewID := map[StateID]StateID{}
	newStIDToRegexID := make(StateIDToRegexID)
	oldStIDToRegexID := nfa.stIDToRegexID
	newq := collection.NewSet[State]()
	qiter := nfa.q.Iterator()
	for qiter.HasNext() {
		oldst := qiter.Next()
		newst := NewState(id)
		newStIDToRegexID.Set(id, oldStIDToRegexID.Get(oldst.GetID()))
		newq.Insert(newst)
		oldToNewID[oldst.GetID()] = id
		id++
	}

	newdelta := make(NFATransition)
	for pair, tos := range nfa.delta {
		oldfrom := pair.First
		b := pair.Second

		newfrom := NewState(oldToNewID[oldfrom.GetID()])
		newtos := collection.NewSet[State]()
		titer := tos.Iterator()
		for titer.HasNext() {
			oldto := titer.Next()
			newto := NewState(oldToNewID[oldto.GetID()])
			newtos.Insert(newto)
		}
		newdelta[collection.NewPair(newfrom, b)] = newtos
	}

	newInitStates := collection.NewSet[State]()
	iiter := nfa.initStates.Iterator()
	for iiter.HasNext() {
		oldst := iiter.Next()
		newst := NewState(oldToNewID[oldst.GetID()])
		newInitStates.Insert(newst)
	}

	newFinStates := collection.NewSet[State]()
	fiter := nfa.finStates.Iterator()
	for fiter.HasNext() {
		oldst := fiter.Next()
		newst := NewState(oldToNewID[oldst.GetID()])
		newFinStates.Insert(newst)
	}

	return NewNFAWithRegexIDMap(newq, newdelta, newInitStates, newFinStates, newStIDToRegexID)
}

func buildStateSet(n int, tos *collection.Set[State]) *StateSet {
	bs := NewStateSet(n)
	titer := tos.Iterator()
	for titer.HasNext() {
		to := titer.Next()
		bs = bs.Insert(to.GetID())
	}
	return bs
}

func (nfa *NFA) SetRegexID(regid RegexID) {
	stIDToRegexID := make(StateIDToRegexID)
	iter := nfa.q.Iterator()
	for iter.HasNext() {
		st := iter.Next()
		if nfa.finStates.Contains(st) {
			stIDToRegexID.Set(st.GetID(), regid)
		} else {
			stIDToRegexID.Set(st.GetID(), nonFinStateRegexID)
		}
	}

	nfa.stIDToRegexID = stIDToRegexID
}
