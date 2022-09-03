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
	q          *collection.Set[State]
	delta      NFATransition
	initStates *collection.Set[State]
	finStates  *collection.Set[State]
	regexID    RegexID
}

func NewNFA(
	q *collection.Set[State],
	delta NFATransition,
	initStates *collection.Set[State],
	finStates *collection.Set[State]) NFA {
	return NFA{
		q: q,
		// sigma:      sigma,
		delta:      delta,
		initStates: initStates,
		finStates:  finStates,
		regexID:    0,
	}
}

func (nfa NFA) Copy() NFA {
	return NewNFA(nfa.q.Copy(), nfa.delta.Copy(), nfa.initStates.Copy(), nfa.finStates.Copy())
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
	n := maxID + 1
	stIDToRegID := make([]RegexID, n)
	qiter := nfa.q.Iterator()
	for qiter.HasNext() {
		st := qiter.Next()
		stIDToRegID[st.GetID()] = st.regexID
	}
	delta := make(ImdNFATransition)
	for pair, tos := range nfa.delta {
		from := pair.First
		b := pair.Second
		delta[collection.NewPair(from.GetID(), b)] = buildStateSet(n, tos)
	}
	initStates := buildStateSet(n, nfa.initStates)
	finStates := buildStateSet(n, nfa.finStates)

	return NewImdNFA(maxID, stIDToRegID, delta, initStates, finStates)
}

func (nfa NFA) relabelStateIDs() NFA {
	nfa = nfa.Copy()
	id := StateID(1)
	oldToNewID := map[StateID]StateID{}
	newq := collection.NewSet[State]()
	qiter := nfa.q.Iterator()
	for qiter.HasNext() {
		oldst := qiter.Next()
		newst := NewState(id)
		newst.SetRegexID(oldst.GetRawRegexID())
		newq.Insert(newst)
		oldToNewID[oldst.GetID()] = id
		id++
	}

	newdelta := make(NFATransition)
	for pair, tos := range nfa.delta {
		oldfrom := pair.First
		b := pair.Second

		newfrom := NewState(oldToNewID[oldfrom.GetID()])
		newfrom.SetRegexID(oldfrom.GetRawRegexID())
		newtos := collection.NewSet[State]()
		titer := tos.Iterator()
		for titer.HasNext() {
			oldto := titer.Next()
			newto := NewState(oldToNewID[oldto.GetID()])
			newto.SetRegexID(oldto.GetRawRegexID())
			newtos.Insert(newto)
		}
		newdelta[collection.NewPair(newfrom, b)] = newtos
	}

	newInitStates := collection.NewSet[State]()
	iiter := nfa.initStates.Iterator()
	for iiter.HasNext() {
		oldst := iiter.Next()
		newst := NewState(oldToNewID[oldst.GetID()])
		newst.SetRegexID(oldst.GetRawRegexID())
		newInitStates.Insert(newst)
	}

	newFinStates := collection.NewSet[State]()
	fiter := nfa.finStates.Iterator()
	for fiter.HasNext() {
		oldst := fiter.Next()
		newst := NewState(oldToNewID[oldst.GetID()])
		newst.SetRegexID(oldst.GetRawRegexID())
		newFinStates.Insert(newst)
	}

	return NewNFA(newq, newdelta, newInitStates, newFinStates)
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

func (nfa *NFA) SetRegexID(id RegexID) {
	nfa2 := nfa.Copy()

	q := collection.NewSet[State]()
	initStates := collection.NewSet[State]()
	finStates := collection.NewSet[State]()
	delta := make(NFATransition)

	qiter := nfa2.q.Iterator()
	for qiter.HasNext() {
		st := qiter.Next()
		if nfa.finStates.Contains(st) {
			st.SetRegexID(id)
		}
		q.Insert(st)
	}
	iiter := nfa2.initStates.Iterator()
	for iiter.HasNext() {
		st := iiter.Next()
		if nfa.finStates.Contains(st) {
			st.SetRegexID(id)
		}
		initStates.Insert(st)
	}
	fiter := nfa2.finStates.Iterator()
	for fiter.HasNext() {
		st := fiter.Next()
		st.SetRegexID(id)
		finStates.Insert(st)
	}
	for pair, sts := range nfa2.delta {
		from := pair.First
		if nfa.finStates.Contains(from) {
			from.SetRegexID(id)
		}
		b := pair.Second
		nss := collection.NewSet[State]()
		siter := sts.Iterator()
		for siter.HasNext() {
			to := siter.Next()
			if nfa.finStates.Contains(to) {
				to.SetRegexID(id)
			}
			nss.Insert(to)
		}
		delta[collection.NewPair(from, b)] = nss
	}

	nfa2 = NewNFA(q, delta, initStates, finStates)
	nfa2.regexID = id

	*nfa = nfa2
}
