package automaton

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

const blackHole = "BH"

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

func (dfa DFA) Copy() DFA {
	return NewDFA(dfa.q.Copy(), dfa.delta.Copy(), dfa.initState, dfa.finStates.Copy())
}

func (dfa DFA) Totalize() DFA {
	dfa = dfa.Copy()
	bhState := NewState(blackHole)
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

// Brzozowski DFA minimization algorithm
func (dfa DFA) Minimize() DFA {
	return dfa.Reverse().ToDFA().Reverse().ToDFA()
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
	ii, si, fi := 0, 0, 0
	for s := range dfa.q {
		n, err := graph.CreateNode(guid.New()) // assign unique node id
		if err != nil {
			return "", err
		}
		if dfa.initState == s {
			e, err := graph.CreateEdge(guid.New(), start, n)
			if err != nil {
				return "", err
			}
			n.SetLabel(fmt.Sprintf("I%v", ii))
			e.SetLabel(string("start"))
			ii++
		}
		if dfa.finStates.Contains(s) {
			n.SetShape(cgraph.DoubleCircleShape)
			n.SetLabel(fmt.Sprintf("F%v", fi))
			fi++
		} else if s.GetLabel() == blackHole {
			n.SetLabel(blackHole)
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v", si))
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
		e, err := graph.CreateEdge(guid.New(), nodes[from], nodes[to])
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
