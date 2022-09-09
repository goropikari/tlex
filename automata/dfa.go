package automata

import "github.com/goropikari/tlex/collection"

type DFATransition struct {
	delta map[StateID]map[Interval]StateID
}

func NewDFATransition() *DFATransition {
	return &DFATransition{
		delta: make(map[StateID]map[Interval]StateID),
	}
}

func (trans *DFATransition) Set(from StateID, intv Interval, to StateID) {
	_, ok := trans.delta[from]
	if !ok {
		trans.delta[from] = map[Interval]StateID{}
	}

	trans.delta[from][intv] = to
}

type DFA struct {
	size        int
	intvs       []Interval
	states      *collection.Set[StateID]
	trans       *DFATransition
	initState   StateID
	finStates   *collection.Set[StateID]
	stIDToRegID StateIDToRegexID
}

// ここで入る intv は dfa.intvs に入っていることを前提としている
func (dfa *DFA) stepIntv(sid StateID, intv Interval) (stateID StateID, nonDeadState bool) {
	retID, ok := dfa.trans.delta[sid][intv]
	return retID, ok
}

func (dfa *DFA) grouping() [][]StateID {
	regIDMap := map[RegexID][]StateID{}
	siter := dfa.states.Iterator()
	for siter.HasNext() {
		sid := siter.Next()
		rid := dfa.stIDToRegID[sid]
		regIDMap[rid] = append(regIDMap[rid], sid)
	}

	grps := make([][]StateID, 0)
	uf := newStateUnionFind(dfa.size)
	for _, sts := range regIDMap {
		for _, sid := range sts[1:] {
			uf.Unite(sts[0], sid)
		}
		grps = append(grps, sts)
	}

	ngrp := len(grps)
	splitted := true
	for splitted {
		splitted = false
		newuf := newStateUnionFind(dfa.size)

		for _, grp := range grps {
			for i, s0 := range grp {
				for _, s1 := range grp[i+1:] {
					same := true
					for _, intv := range dfa.intvs {
						ns0, ok1 := dfa.stepIntv(s0, intv)
						ns1, ok2 := dfa.stepIntv(s1, intv)

						if ok1 != ok2 {
							same = false
							break
						} else if !ok1 {
							continue
						}

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
		for stID := StateID(0); stID < StateID(dfa.size); stID++ {
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

	return grps
}

func (dfa *DFA) LexerMinimize() *DFA {
	grps := dfa.grouping()

	uf := newStateUnionFind(dfa.size)
	for _, grp := range grps {
		n := len(grp)
		if n == 1 {
			continue
		}

		for i := 1; i < n; i++ {
			uf.Unite(grp[0], grp[i])
		}
	}

	states := collection.NewSet[StateID]()
	stIDToRegID := NewStateIDToRegexID()
	siter := dfa.states.Iterator()
	for siter.HasNext() {
		sid := siter.Next()
		leaderID := uf.Find(sid)
		states.Insert(leaderID)
		stIDToRegID.Set(leaderID, dfa.stIDToRegID.Get(leaderID))
	}

	initState := uf.Find(dfa.initState)

	trans := NewDFATransition()
	for from, mp := range dfa.trans.delta {
		for intv, to := range mp {
			from = uf.Find(from)
			to = uf.Find(to)
			trans.Set(from, intv, to)
		}
	}

	finStates := collection.NewSet[StateID]()
	fiter := dfa.finStates.Iterator()
	for fiter.HasNext() {
		sid := fiter.Next()
		finStates.Insert(uf.Find(sid))
	}

	return &DFA{
		size:        states.Size(),
		intvs:       dfa.intvs,
		states:      states,
		trans:       trans,
		initState:   initState,
		finStates:   finStates,
		stIDToRegID: stIDToRegID,
	}
}
