package automata

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/utils/guid"
	"golang.org/x/exp/slices"
)

const blackHoleStateID = 0

type DFATransition map[collection.Tuple[State, rune]]State

func (t DFATransition) Copy() DFATransition {
	delta := make(DFATransition)
	for k, v := range t {
		delta[k] = v
	}

	return delta
}

type DFA struct {
	q         collection.Set[State]
	delta     DFATransition
	initState State
	finStates collection.Set[State]
}

func NewDFA(q collection.Set[State], delta DFATransition, initState State, finStates collection.Set[State]) DFA {
	return DFA{
		q:         q,
		delta:     delta,
		initState: initState,
		finStates: finStates,
	}
}

func (dfa DFA) GetStates() collection.Set[State] {
	return dfa.q
}

func (dfa DFA) GetInitState() State {
	return dfa.initState
}

func (dfa DFA) GetFinStates() collection.Set[State] {
	return dfa.finStates
}

func (dfa DFA) GetTransitionTable() DFATransition {
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

func (dfa DFA) Accept(s string) (TokenID, bool) {
	currSt := dfa.initState

	for _, ru := range []rune(s) {
		currSt = dfa.Step(currSt, ru)
		if currSt.GetID() == blackHoleStateID {
			return 0, false
		}
		if (currSt == State{}) {
			return 0, false
		}
	}

	return currSt.GetTokenID(), dfa.finStates.Contains(currSt)
}

func (dfa DFA) Step(st State, ru rune) State {
	pair := collection.NewTuple(st, ru)
	return dfa.delta[pair]
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
		for st := range dfa.q {
			tu := collection.NewTuple(st, ru)
			if _, ok := dfa.delta[tu]; !ok {
				changed = true
				delta[tu] = bhState
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
	for pair, ns := range dfa.delta {
		from := pair.First
		ru := pair.Second
		tu := collection.NewTuple(ns, ru)
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

	for pair, to := range dfa.delta {
		if to.GetID() == blackHoleStateID {
			delete(dfa.delta, pair)
		}
	}

	return dfa
}

type stateGroup struct {
	states collection.Set[State]
}

func NewGroup(states collection.Set[State]) *stateGroup {
	return &stateGroup{states: states}
}

func (g *stateGroup) size() int {
	return len(g.states)
}

func (g *stateGroup) slice() []State {
	sts := make([]State, 0)
	for st := range g.states {
		sts = append(sts, st)
	}

	return sts
}

type stateUnionFind struct {
	stToID map[State]int
	idToSt map[int]State
	uf     *collection.UnionFind
}

func newStateUnionFind(states []State) *stateUnionFind {
	stToID := make(map[State]int)
	idToSt := make(map[int]State)
	id := 0
	for _, st := range states {
		stToID[st] = id
		idToSt[id] = st
		id++
	}

	return &stateUnionFind{
		stToID: stToID,
		idToSt: idToSt,
		uf:     collection.NewUnionFind(len(states)),
	}
}

func (uf *stateUnionFind) unite(x, y State) bool {
	xid := uf.stToID[x]
	yid := uf.stToID[y]
	return uf.uf.Unite(xid, yid)
}

func (uf *stateUnionFind) find(x State) State {
	id := uf.stToID[x]
	return uf.idToSt[uf.uf.Find(id)]
}

// state minimization for lexical analyzer
// Compilers: Principles, Techniques, and Tools, 2ed ed.,  ISBN 9780321486813 (Dragon book)
// p.181 Algorithm 3.39
// p.184 3.9.7 State Minimization in Lexical Analyzers
func (dfa DFA) grouping() []*stateGroup {
	states := dfa.q.Slice()

	stateSets := make(map[TokenID]collection.Set[State])
	for st := range dfa.q {
		if _, ok := stateSets[st.GetTokenID()]; ok {
			stateSets[st.GetTokenID()].Insert(st)
		} else {
			stateSets[st.GetTokenID()] = collection.NewSet[State]().Insert(st)
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
						ns0 := dfa.delta[collection.NewTuple(s0, ru)]
						ns1 := dfa.delta[collection.NewTuple(s1, ru)]
						// If ns0 and ns1 belong to different groups, s0 and s1 belong to other groups.
						// Then current group is split.
						if oldStUF.find(ns0) != oldStUF.find(ns1) {
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

		newStateSets := make(map[State]collection.Set[State])
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
	for st := range dfa.q {
		q.Insert(uf.find(st))
	}

	initState := uf.find(dfa.initState)

	delta := make(DFATransition)
	for pair, ns := range dfa.delta {
		from := uf.find(pair.First)
		ru := pair.Second
		ns = uf.find(ns)
		delta[collection.NewTuple(from, ru)] = ns
	}

	finStates := collection.NewSet[State]()
	for st := range dfa.finStates {
		finStates.Insert(uf.find(st))
	}

	return NewDFA(q, delta, initState, finStates)
}

func (dfa DFA) ToDot() (string, error) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()
	graph.SetRankDir("LR") // 図を横長にする

	start, err := graph.CreateNode("start")
	if err != nil {
		return "", err
	}
	start.SetShape(cgraph.PointShape)
	nodes := make(map[State]*cgraph.Node)
	si, fi := 0, 0
	for s := range dfa.q {
		n, err := graph.CreateNode(fmt.Sprintf("%v", guid.New())) // assign unique node id
		if err != nil {
			return "", err
		}
		if dfa.initState == s {
			e, err := graph.CreateEdge(fmt.Sprintf("%v", guid.New()), start, n)
			if err != nil {
				return "", err
			}
			e.SetLabel(string("start"))
		}
		if dfa.finStates.Contains(s) {
			n.SetShape(cgraph.DoubleCircleShape)
			n.SetLabel(fmt.Sprintf("F%v_%v", fi, toStateTokenID(s.GetTokenID())))
			fi++
		} else if s.GetID() == blackHoleStateID {
			// n.SetLabel(blackHole)
			n.SetLabel("BH")
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v_%v", si, toStateTokenID(s.GetTokenID())))
			si++
		}
		nodes[s] = n
	}

	// add edge labels
	edges := make(map[collection.Tuple[State, State]][]string)
	for st, to := range dfa.delta {
		from := st.First
		symbol := charLabel(string(st.Second))
		edges[collection.NewTuple(from, to)] = append(edges[collection.NewTuple(from, to)], symbol)
	}
	for edge, labels := range edges {
		from, to := edge.First, edge.Second
		e, err := graph.CreateEdge(fmt.Sprintf("%v", guid.New()), nodes[from], nodes[to])
		if err != nil {
			return "", err
		}
		slices.Sort(labels)
		e.SetLabel(strings.Join(labels, "\n"))
	}

	var buf bytes.Buffer
	g.Render(graph, "dot", &buf)

	return buf.String(), nil
}

func toStateTokenID(id TokenID) TokenID {
	if id == nonFinStateTokenID {
		return 0
	}

	return id
}
