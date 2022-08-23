package automaton

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/google/uuid"
	"github.com/goropikari/golex/collection"
	"golang.org/x/exp/slices"
)

const SupportedChars = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ \t\n\r"

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
			if _, ok := nfa.delta[collection.NewTuple(from, epsilon)]; ok {
				nfa.delta[collection.NewTuple(from, epsilon)].Insert(to)
			} else {
				nfa.delta[collection.NewTuple(from, epsilon)] = collection.NewSet[State]().Insert(to)
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

	nfa.delta[collection.NewTuple(startFinState, epsilon)] = nfa.initStates

	for from := range nfa.finStates {
		nfa.delta[collection.NewTuple(from, epsilon)] = initStates
	}

	// return NewNFA(nfa.q, nfa.sigma, nfa.delta, initStates, initStates)
	return NewNFA(nfa.q, nfa.delta, initStates, initStates)
}

func (nfa NFA) ToDFA() DFA {
	que := list.New()
	memo := collection.NewSet[string]()
	initStates := nfa.eClosureSet(nfa.initStates)
	finStates := nfa.eClosureSet(nfa.finStates)
	que.PushBack(initStates)
	memo.Insert(labelConcat(initStates))

	// dfaInitStates := collection.NewSet[State]().Insert(NewState(labelConcat(initStates)))
	dfaFinStates := collection.NewSet[State]()
	if len(initStates.Intersection(finStates)) > 0 {
		dfaFinStates.Insert(NewState(labelConcat(finStates)))
	}
	dfaDelta := make(Transition)

	for que.Len() > 0 {
		top := que.Front()
		que.Remove(top)
		froms := top.Value.(collection.Set[State])
		fromLabel := labelConcat(froms)

		for _, ru := range SupportedChars {
			tos := collection.NewSet[State]()
			for from := range froms {
				if nx, ok := nfa.delta[collection.NewTuple(from, ru)]; ok {
					tos = tos.Union(nfa.eClosureSet(nx))
				}
			}
			if len(tos) == 0 {
				continue
			}
			toLabel := labelConcat(tos)
			if len(tos.Intersection(finStates)) > 0 {
				dfaFinStates.Insert(NewState(toLabel))
			}
			dfaDelta[collection.NewTuple(NewState(fromLabel), ru)] = collection.NewSet[State]().Insert(NewState(toLabel))
			if memo.Contains(toLabel) {
				continue
			}
			memo.Insert(toLabel)
			que.PushBack(tos)
		}
	}

	q := collection.NewSet[State]()
	for label := range memo {
		q.Insert(NewState(label))
	}

	return NewDFA(q, dfaDelta, NewState(labelConcat(initStates)), dfaFinStates)
}

func (nfa NFA) eClosure(st State) collection.Set[State] {
	que := list.New()
	que.PushBack(st)
	visited := collection.NewSet[State]().Insert(st)

	closure := visited.Copy()
	for que.Len() > 0 {
		front := que.Front()
		que.Remove(front)
		top := front.Value.(State)

		if nxs, ok := nfa.delta[collection.NewTuple(top, epsilon)]; ok {
			closure = closure.Union(nxs)

			for nx := range nxs {
				if !visited.Contains(nx) {
					visited.Insert(nx)
					que.PushBack(nx)
				}
			}
		}
	}

	return closure
}

func (nfa NFA) eClosureSet(sts collection.Set[State]) collection.Set[State] {
	closure := collection.NewSet[State]()
	for st := range sts {
		closure = closure.Union(nfa.eClosure(st))
	}
	return closure
}

func labelConcat(set collection.Set[State]) string {
	s := make([]string, 0, len(set))
	for v := range set {
		s = append(s, v.GetLabel())
	}
	slices.Sort(s)
	return strings.Join(s, "_")
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
		from := st.First
		symbol := string(st.Second)
		for to := range qs {
			e, err := graph.CreateEdge(charLabel(symbol), nodes[from], nodes[to])
			if err != nil {
				return "", err
			}
			e.SetLabel(charLabel(symbol))
		}
	}

	var buf bytes.Buffer
	g.Render(graph, "dot", &buf)

	return buf.String(), nil
}
