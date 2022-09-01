package automata

import (
	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/utils/guid"
)

type NFATransition map[collection.Tuple[State, rune]]collection.Set[State]

func (t NFATransition) Copy() NFATransition {
	delta := make(NFATransition)
	for k, v := range t {
		delta[k] = v.Copy()
	}

	return delta
}

type NFA struct {
	q collection.Set[State]
	// sigma      collection.Set[rune]
	delta      NFATransition
	initStates collection.Set[State]
	finStates  collection.Set[State]
	regexID    RegexID
}

func NewNFA(
	q collection.Set[State],
	// sigma collection.Set[rune],
	delta NFATransition,
	initStates collection.Set[State],
	finStates collection.Set[State]) NFA {
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

	for st := range other.q {
		nfa.q.Insert(st)
	}

	for tr, ss := range other.delta {
		nfa.delta[tr] = ss
	}

	for from := range nfa.finStates {
		for to := range other.initStates {
			if _, ok := nfa.delta[collection.NewTuple(from, epsilon)]; ok {
				nfa.delta[collection.NewTuple(from, epsilon)].Insert(to)
			} else {
				nfa.delta[collection.NewTuple(from, epsilon)] = collection.NewSet[State]().Insert(to)
			}
		}
	}

	return NewNFA(nfa.q, nfa.delta, nfa.initStates, other.finStates)
}

func (nfa NFA) Sum(other NFA) NFA {
	nfa = nfa.Copy()
	other = other.Copy()

	for st := range other.q {
		nfa.q.Insert(st)
	}

	for tr, ss := range other.delta {
		nfa.delta[tr] = ss
	}

	for st := range other.initStates {
		nfa.initStates.Insert(st)
	}

	for st := range other.finStates {
		nfa.finStates.Insert(st)
	}

	return NewNFA(nfa.q, nfa.delta, nfa.initStates, nfa.finStates)
}

func (nfa NFA) Star() NFA {
	nfa = nfa.Copy()

	startFinState := NewState(StateID(guid.New()))
	initStates := collection.NewSet[State]().Insert(startFinState)

	nfa.q.Insert(startFinState)

	nfa.delta[collection.NewTuple(startFinState, epsilon)] = nfa.initStates

	for from := range nfa.finStates {
		pair := collection.NewTuple(from, epsilon)
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
	maxID := len(nfa.q)
	n := maxID + 1
	stIDToRegID := make([]RegexID, n)
	for st := range nfa.q {
		stIDToRegID[int(st.GetID())] = st.regexID
	}
	delta := make(ImdNFATransition)
	for pair, tos := range nfa.delta {
		from := pair.First
		ru := pair.Second
		delta[collection.NewTuple(from.GetID(), ru)] = buildStateSet(n, tos)
	}
	initStates := buildStateSet(n, nfa.initStates)
	finStates := buildStateSet(n, nfa.finStates)

	return ImdNFA{
		maxID:       maxID,
		stIDToRegID: stIDToRegID,
		delta:       delta,
		initStates:  initStates,
		finStates:   finStates,
	}
}

func (nfa NFA) relabelStateIDs() NFA {
	nfa = nfa.Copy()
	id := StateID(1)
	oldToNewID := map[StateID]StateID{}
	newq := collection.NewSet[State]()
	for oldst := range nfa.q {
		newst := NewState(id)
		newst.SetRegexID(oldst.GetRawRegexID())
		newq.Insert(newst)
		oldToNewID[oldst.GetID()] = id
		id++
	}

	newdelta := make(NFATransition)
	for pair, tos := range nfa.delta {
		oldfrom := pair.First
		ru := pair.Second

		newfrom := NewState(oldToNewID[oldfrom.GetID()])
		newfrom.SetRegexID(oldfrom.GetRawRegexID())
		newtos := collection.NewSet[State]()
		for oldto := range tos {
			newto := NewState(oldToNewID[oldto.GetID()])
			newto.SetRegexID(oldto.GetRawRegexID())
			newtos.Insert(newto)
		}
		newdelta[collection.NewTuple(newfrom, ru)] = newtos
	}

	newInitStates := collection.NewSet[State]()
	for oldst := range nfa.initStates {
		newst := NewState(oldToNewID[oldst.GetID()])
		newst.SetRegexID(oldst.GetRawRegexID())
		newInitStates.Insert(newst)
	}

	newFinStates := collection.NewSet[State]()
	for oldst := range nfa.finStates {
		newst := NewState(oldToNewID[oldst.GetID()])
		newst.SetRegexID(oldst.GetRawRegexID())
		newFinStates.Insert(newst)
	}

	return NewNFA(newq, newdelta, newInitStates, newFinStates)
}

func buildStateSet(n int, tos collection.Set[State]) *StateSet {
	bs := NewStateSet(n)
	for to := range tos {
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

	for st := range nfa2.q {
		if nfa.finStates.Contains(st) {
			st.SetRegexID(id)
		}
		q.Insert(st)
	}
	for st := range nfa2.initStates {
		if nfa.finStates.Contains(st) {
			st.SetRegexID(id)
		}
		initStates.Insert(st)
	}
	for st := range nfa2.finStates {
		st.SetRegexID(id)
		finStates.Insert(st)
	}
	for pair, sts := range nfa2.delta {
		from := pair.First
		if nfa.finStates.Contains(from) {
			from.SetRegexID(id)
		}
		ru := pair.Second
		nss := collection.NewSet[State]()
		for to := range sts {
			if nfa.finStates.Contains(to) {
				to.SetRegexID(id)
			}
			nss.Insert(to)
		}
		delta[collection.NewTuple(from, ru)] = nss
	}

	nfa2 = NewNFA(q, delta, initStates, finStates)
	nfa2.regexID = id

	*nfa = nfa2
}
