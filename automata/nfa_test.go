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
		automata.NewEpsilonTransition(),
		automata.NewNFATransition().Set(id0, automata.NewInterval(65, 65), id1),
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
		automata.NewEpsilonTransition(),
		automata.NewNFATransition().
			Set(id2, automata.NewInterval(65, 65), id3).
			Set(id4, automata.NewInterval(66, 66), id5),
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
		automata.NewEpsilonTransition(),
		automata.NewNFATransition().
			Set(id7, automata.NewInterval(65, 65), id8).
			Set(id9, automata.NewInterval(66, 66), id10).
			Set(id12, automata.NewInterval(66, 66), id13),
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
		automata.NewEpsilonTransition(),
		automata.NewNFATransition().
			Set(id0, automata.NewInterval(12354, 12358), id1),
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
		automata.NewEpsilonTransition().Set(id3, id4),
		automata.NewNFATransition().
			Set(id2, automata.NewInterval(12354, 12356), id3).
			Set(id4, automata.NewInterval(66, 70), id5),
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
		automata.NewEpsilonTransition().Set(id6, id7).Set(id6, id9).Set(id8, id6).Set(id10, id11).Set(id11, id12).Set(id13, id11),
		automata.NewNFATransition().
			Set(id7, automata.NewInterval(66, 68), id8).
			Set(id9, automata.NewInterval(66, 66), id10).
			Set(id12, automata.NewInterval(66, 66), id13),
		collection.NewSet[automata.StateID]().Insert(id6),
		collection.NewSet[automata.StateID]().Insert(id11),
	)
	nfa3.SetRegexID(3)

	nfa := nfa1.Sum(nfa2).Sum(nfa3).ToImdNFA().ToDFA().LexerMinimize()
	fmt.Println(nfa.ToDot())
}
