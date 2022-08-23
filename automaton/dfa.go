package automaton

import (
	"bytes"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/google/uuid"
	"github.com/goropikari/golex/collection"
)

type DFA struct {
	q         collection.Set[State]
	delta     Transition
	initState State
	finStates collection.Set[State]
}

func NewDFA(q collection.Set[State], delta Transition, initState State, finStates collection.Set[State]) DFA {
	return DFA{
		q:         q,
		delta:     delta,
		initState: initState,
		finStates: finStates,
	}
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
		n, err := graph.CreateNode(uuid.New().String()) // assign unique node id
		if err != nil {
			return "", err
		}
		if dfa.initState == s {
			e, err := graph.CreateEdge(uuid.New().String(), start, n)
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
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v", si))
			si++
		}
		nodes[s] = n
	}

	for st, qs := range dfa.delta {
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
