package automata_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/tlex/automata"
	"github.com/goropikari/tlex/collection"
)

func TestNFA(t *testing.T) {
	automata.NewStateID()
	automata.NewStateID()
	automata.NewStateID()

	id0 := automata.NewStateID()
	id1 := automata.NewStateID()
	nfa1 := automata.NewNFA(
		collection.NewSet[automata.StateID]().Insert(id0).Insert(id1),
		automata.NewEpsilonTransition(make(map[automata.StateID]*collection.Set[automata.StateID])),
		automata.NewTransition(
			map[automata.StateID]map[automata.Interval]*collection.Set[automata.StateID]{
				id0: {
					automata.NewInterval(65, 65): collection.NewSet[automata.StateID]().Insert(id1),
				},
			},
		),
		collection.NewSet[automata.StateID]().Insert(id0),
		collection.NewSet[automata.StateID]().Insert(id1),
	)
	nfa1.SetRegexID(1)

	id2 := automata.NewStateID()
	id3 := automata.NewStateID()
	id4 := automata.NewStateID()
	id5 := automata.NewStateID()
	nfa2 := automata.NewNFA(
		collection.NewSet[automata.StateID]().Insert(id2).Insert(id3).Insert(id4).Insert(id5),
		automata.NewEpsilonTransition(map[automata.StateID]*collection.Set[automata.StateID]{
			id3: collection.NewSet[automata.StateID]().Insert(id4),
		}),
		automata.NewTransition(
			map[automata.StateID]map[automata.Interval]*collection.Set[automata.StateID]{
				id2: {
					automata.NewInterval(65, 65): collection.NewSet[automata.StateID]().Insert(id3),
				},
				id4: {
					automata.NewInterval(66, 66): collection.NewSet[automata.StateID]().Insert(id5),
				},
			},
		),
		collection.NewSet[automata.StateID]().Insert(id2),
		collection.NewSet[automata.StateID]().Insert(id5),
	)
	nfa2.SetRegexID(2)

	id6 := automata.NewStateID()
	id7 := automata.NewStateID()
	id8 := automata.NewStateID()
	id9 := automata.NewStateID()
	id10 := automata.NewStateID()
	id11 := automata.NewStateID()
	id12 := automata.NewStateID()
	id13 := automata.NewStateID()
	nfa3 := automata.NewNFA(
		collection.NewSet[automata.StateID]().Insert(id6).Insert(id7).Insert(id8).Insert(id9).Insert(id10).Insert(id11).Insert(id12).Insert(id13),
		automata.NewEpsilonTransition(map[automata.StateID]*collection.Set[automata.StateID]{
			id6:  collection.NewSet[automata.StateID]().Insert(id7).Insert(id9),
			id8:  collection.NewSet[automata.StateID]().Insert(id6),
			id10: collection.NewSet[automata.StateID]().Insert(id11),
			id11: collection.NewSet[automata.StateID]().Insert(id12),
			id13: collection.NewSet[automata.StateID]().Insert(id11),
		}),
		automata.NewTransition(
			map[automata.StateID]map[automata.Interval]*collection.Set[automata.StateID]{
				id7: {
					automata.NewInterval(65, 65): collection.NewSet[automata.StateID]().Insert(id8),
				},
				id9: {
					automata.NewInterval(66, 66): collection.NewSet[automata.StateID]().Insert(id10),
				},
				id12: {
					automata.NewInterval(66, 66): collection.NewSet[automata.StateID]().Insert(id13),
				},
			},
		),
		collection.NewSet[automata.StateID]().Insert(id6),
		collection.NewSet[automata.StateID]().Insert(id11),
	)
	nfa3.SetRegexID(3)

	// a|ab|a*bb*
	nfa := nfa1.Sum(nfa2).Sum(nfa3).ToImdNFA().ToDFA().LexerMinimize()
	nfa.ToDot()
}

func TestNFA2(t *testing.T) {
	automata.NewStateID()
	automata.NewStateID()
	automata.NewStateID()

	id0 := automata.NewStateID()
	id1 := automata.NewStateID()
	nfa1 := automata.NewNFA(
		collection.NewSet[automata.StateID]().Insert(id0).Insert(id1),
		automata.NewEpsilonTransition(make(map[automata.StateID]*collection.Set[automata.StateID])),
		automata.NewTransition(
			map[automata.StateID]map[automata.Interval]*collection.Set[automata.StateID]{
				id0: {
					automata.NewInterval(12354, 12358): collection.NewSet[automata.StateID]().Insert(id1),
				},
			},
		),
		collection.NewSet[automata.StateID]().Insert(id0),
		collection.NewSet[automata.StateID]().Insert(id1),
	)
	nfa1.SetRegexID(1)

	id2 := automata.NewStateID()
	id3 := automata.NewStateID()
	id4 := automata.NewStateID()
	id5 := automata.NewStateID()
	nfa2 := automata.NewNFA(
		collection.NewSet[automata.StateID]().Insert(id2).Insert(id3).Insert(id4).Insert(id5),
		automata.NewEpsilonTransition(map[automata.StateID]*collection.Set[automata.StateID]{
			id3: collection.NewSet[automata.StateID]().Insert(id4),
		}),
		automata.NewTransition(
			map[automata.StateID]map[automata.Interval]*collection.Set[automata.StateID]{
				id2: {
					automata.NewInterval(12354, 12356): collection.NewSet[automata.StateID]().Insert(id3),
				},
				id4: {
					automata.NewInterval(66, 70): collection.NewSet[automata.StateID]().Insert(id5),
				},
			},
		),
		collection.NewSet[automata.StateID]().Insert(id2),
		collection.NewSet[automata.StateID]().Insert(id5),
	)
	nfa2.SetRegexID(2)

	id6 := automata.NewStateID()
	id7 := automata.NewStateID()
	id8 := automata.NewStateID()
	id9 := automata.NewStateID()
	id10 := automata.NewStateID()
	id11 := automata.NewStateID()
	id12 := automata.NewStateID()
	id13 := automata.NewStateID()
	nfa3 := automata.NewNFA(
		collection.NewSet[automata.StateID]().Insert(id6).Insert(id7).Insert(id8).Insert(id9).Insert(id10).Insert(id11).Insert(id12).Insert(id13),
		automata.NewEpsilonTransition(map[automata.StateID]*collection.Set[automata.StateID]{
			id6:  collection.NewSet[automata.StateID]().Insert(id7).Insert(id9),
			id8:  collection.NewSet[automata.StateID]().Insert(id6),
			id10: collection.NewSet[automata.StateID]().Insert(id11),
			id11: collection.NewSet[automata.StateID]().Insert(id12),
			id13: collection.NewSet[automata.StateID]().Insert(id11),
		}),
		automata.NewTransition(
			map[automata.StateID]map[automata.Interval]*collection.Set[automata.StateID]{
				id7: {
					automata.NewInterval(66, 68): collection.NewSet[automata.StateID]().Insert(id8),
				},
				id9: {
					automata.NewInterval(66, 66): collection.NewSet[automata.StateID]().Insert(id10),
				},
				id12: {
					automata.NewInterval(66, 66): collection.NewSet[automata.StateID]().Insert(id13),
				},
			},
		),
		collection.NewSet[automata.StateID]().Insert(id6),
		collection.NewSet[automata.StateID]().Insert(id11),
	)
	nfa3.SetRegexID(3)

	// a|ab|a*bb*
	nfa := nfa1.Sum(nfa2).Sum(nfa3).ToImdNFA().ToDFA().LexerMinimize()
	fmt.Println(nfa.ToDot())
}
