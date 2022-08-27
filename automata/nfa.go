package automata

import (
	"bytes"
	"container/list"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
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
	tokenID    TokenID
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
		tokenID:    0,
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

	startFinState := NewState(guid.New())
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

func (nfa NFA) ToDFA() DFA {
	que := list.New()
	initStates := nfa.eClosureSet(nfa.initStates)
	que.PushBack(initStates)
	finStates := nfa.finStates
	dfaInitStates := NewStateSet(initStates)
	memo := collection.NewSet[State]().Insert(dfaInitStates)

	dfaFinStates := collection.NewSet[State]()
	if len(initStates.Intersection(finStates)) > 0 {
		dfaFinStates.Insert(dfaInitStates)
	}
	dfaDelta := make(DFATransition)

	for que.Len() > 0 {
		top := que.Front()
		que.Remove(top)
		froms := top.Value.(collection.Set[State])

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
			to := NewStateSet(tos)
			if len(tos.Intersection(finStates)) > 0 {
				dfaFinStates.Insert(to)
			}
			dfaDelta[collection.NewTuple(NewStateSet(froms), ru)] = to
			if memo.Contains(to) {
				continue
			}
			memo.Insert(to)
			que.PushBack(tos)
		}
	}

	q := collection.NewSet[State]()
	for st := range memo {
		q.Insert(st)
	}

	return NewDFA(q, dfaDelta, dfaInitStates, dfaFinStates)
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

func (nfa *NFA) SetTokenID(id TokenID) {
	nfa2 := nfa.Copy()

	q := collection.NewSet[State]()
	initStates := collection.NewSet[State]()
	finStates := collection.NewSet[State]()
	delta := make(NFATransition)

	for st := range nfa2.q {
		if nfa.finStates.Contains(st) {
			st.SetTokenID(id)
		}
		q.Insert(st)
	}
	for st := range nfa2.initStates {
		if nfa.finStates.Contains(st) {
			st.SetTokenID(id)
		}
		initStates.Insert(st)
	}
	for st := range nfa2.finStates {
		st.SetTokenID(id)
		finStates.Insert(st)
	}
	for pair, sts := range nfa2.delta {
		from := pair.First
		if nfa.finStates.Contains(from) {
			from.SetTokenID(id)
		}
		ru := pair.Second
		nss := collection.NewSet[State]()
		for to := range sts {
			if nfa.finStates.Contains(to) {
				to.SetTokenID(id)
			}
			nss.Insert(to)
		}
		delta[collection.NewTuple(from, ru)] = nss
	}

	nfa2 = NewNFA(q, delta, initStates, finStates)
	nfa2.tokenID = id

	*nfa = nfa2
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
		n, err := graph.CreateNode(guid.New()) // assign unique node id
		if err != nil {
			return "", err
		}
		if nfa.initStates.Contains(s) {
			e, err := graph.CreateEdge(guid.New(), start, n)
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
