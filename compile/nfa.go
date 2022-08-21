package compile

import (
	"bytes"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/google/uuid"
	"github.com/goropikari/golex/collection"
)

const epsilon = 'ε'

type Transition map[Tuple[State, rune]]collection.Set[State]

func (t Transition) Copy() Transition {
	delta := make(Transition)
	for k, v := range t {
		delta[k] = v.Copy()
	}

	return delta
}

type State struct {
	label string
}

func NewState(label string) State {
	return State{label: label}
}

type NFA struct {
	q collection.Set[State]
	// sigma      collection.Set[rune]
	delta      Transition
	initStates collection.Set[State]
	finStates  collection.Set[State]
}

func NewNFA(
	q collection.Set[State],
	// sigma collection.Set[rune],
	delta Transition,
	initStates collection.Set[State],
	finStates collection.Set[State]) NFA {
	return NFA{
		q: q,
		// sigma:      sigma,
		delta:      delta,
		initStates: initStates,
		finStates:  finStates,
	}
}

func (nfa NFA) Copy() NFA {
	// return NewNFA(nfa.q.Copy(), nfa.sigma.Copy(), nfa.delta.Copy(), nfa.initStates.Copy(), nfa.finStates.Copy())
	return NewNFA(nfa.q.Copy(), nfa.delta.Copy(), nfa.initStates.Copy(), nfa.finStates.Copy())
}

func (nfa NFA) Concat(other NFA) NFA {
	nfa = nfa.Copy()
	other = other.Copy()

	for st := range other.q {
		nfa.q.Insert(st)
	}

	// for r := range other.sigma {
	// 	nfa.sigma.Insert(r)
	// }

	for tr, ss := range other.delta {
		nfa.delta[tr] = ss
	}

	for from := range nfa.finStates {
		for to := range other.initStates {
			if _, ok := nfa.delta[NewTuple(from, epsilon)]; ok {
				nfa.delta[NewTuple(from, epsilon)].Insert(to)
			} else {
				nfa.delta[NewTuple(from, epsilon)] = collection.NewSet[State]().Insert(to)
			}
		}
	}

	// return NewNFA(nfa.q, nfa.sigma, nfa.delta, nfa.initStates, other.finStates)
	return NewNFA(nfa.q, nfa.delta, nfa.initStates, other.finStates)
}

func (nfa NFA) Sum(other NFA) NFA {
	nfa = nfa.Copy()
	other = other.Copy()

	for st := range other.q {
		nfa.q.Insert(st)
	}

	// for r := range other.sigma {
	// 	nfa.sigma.Insert(r)
	// }

	for tr, ss := range other.delta {
		nfa.delta[tr] = ss
	}

	for st := range other.initStates {
		nfa.initStates.Insert(st)
	}

	for st := range other.finStates {
		nfa.finStates.Insert(st)
	}

	// return NewNFA(nfa.q, nfa.sigma, nfa.delta, nfa.initStates, nfa.finStates)
	return NewNFA(nfa.q, nfa.delta, nfa.initStates, nfa.finStates)
}

func (nfa NFA) Star() NFA {
	nfa = nfa.Copy()

	startFinState := NewState(uuid.New().String())
	initStates := collection.NewSet[State]().Insert(startFinState)

	nfa.q.Insert(startFinState)

	nfa.delta[NewTuple(startFinState, epsilon)] = nfa.initStates

	for from := range nfa.finStates {
		nfa.delta[NewTuple(from, epsilon)] = initStates
	}

	// return NewNFA(nfa.q, nfa.sigma, nfa.delta, initStates, initStates)
	return NewNFA(nfa.q, nfa.delta, initStates, initStates)
}

func (nfa NFA) ToDot() (string, error) {
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
	nodes := make(map[State]*cgraph.Node)
	ii, si, fi := 0, 0, 0
	for s := range nfa.q {
		n, err := graph.CreateNode(uuid.New().String()) // assign unique node id
		if err != nil {
			return "", err
		}
		if nfa.initStates.Contains(s) {
			e, err := graph.CreateEdge(uuid.New().String(), start, n)
			if err != nil {
				return "", err
			}
			n.SetLabel(fmt.Sprintf("I%v", ii))
			e.SetLabel(string(epsilon))
			ii++
		}
		if nfa.finStates.Contains(s) {
			n.SetShape(cgraph.DoubleCircleShape)
			n.SetLabel(fmt.Sprintf("F%v", fi))
			fi++
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v", si))
			si++
		}
		nodes[s] = n
	}

	for st, qs := range nfa.delta {
		from := st.first
		label := string(st.second)
		for to := range qs {
			e, err := graph.CreateEdge(label, nodes[from], nodes[to])
			if err != nil {
				return "", err
			}
			e.SetLabel(label)
		}
	}

	var buf bytes.Buffer
	g.Render(graph, "dot", &buf)

	return buf.String(), nil
}
