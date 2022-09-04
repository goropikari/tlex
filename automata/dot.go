package automata

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/goropikari/tlex/collection"
	"github.com/goropikari/tlex/utils/guid"
	"golang.org/x/exp/slices"
)

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
	qiter := nfa.q.Iterator()
	for qiter.HasNext() {
		s := qiter.Next()
		n, err := graph.CreateNode(fmt.Sprintf("%v", guid.New())) // assign unique node id
		if err != nil {
			return "", err
		}
		if nfa.initStates.Contains(s) {
			e, err := graph.CreateEdge(fmt.Sprintf("%v", guid.New()), start, n)
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
		qsiter := qs.Iterator()
		for qsiter.HasNext() {
			to := qsiter.Next()
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

func (nfa ImdNFA) ToDot() (string, error) {
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
	for id := 1; id <= nfa.maxID; id++ {
		sid := StateID(id)
		n, err := graph.CreateNode(fmt.Sprintf("%v", guid.New())) // assign unique node id
		if err != nil {
			return "", err
		}
		if nfa.initStates.Contains(sid) {
			e, err := graph.CreateEdge(fmt.Sprintf("%v", guid.New()), start, n)
			if err != nil {
				return "", err
			}
			n.SetLabel(fmt.Sprintf("I%v", ii))
			e.SetLabel(string(epsilon))
			ii++
		}
		if nfa.finStates.Contains(sid) {
			n.SetShape(cgraph.DoubleCircleShape)
			n.SetLabel(fmt.Sprintf("F%v", fi))
			fi++
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v", si))
			si++
		}
		st := NewState(StateID(id))
		nodes[st] = n
	}

	for st, qs := range nfa.delta {
		from := st.First
		symbol := string(st.Second)
		fromst := NewState(from)
		for id := 1; id <= nfa.maxID; id++ {
			sid := StateID(id)
			if !qs.Contains(sid) {
				continue
			}
			tost := NewState(sid)
			e, err := graph.CreateEdge(charLabel(symbol), nodes[fromst], nodes[tost])
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
	qiter := dfa.q.Iterator()
	for qiter.HasNext() {
		s := qiter.Next()
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
			n.SetLabel(fmt.Sprintf("F%v_%v", fi, toStateRegexID(dfa.GetRegexID(s))))
			fi++
		} else if s.GetID() == blackHoleStateID {
			n.SetLabel("BH")
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v_%v", si, toStateRegexID(dfa.GetRegexID(s))))
			si++
		}
		nodes[s] = n
	}

	// add edge labels
	edges := make(map[collection.Pair[State, State]][]string)
	iter := dfa.delta.Iterator()
	for iter.HasNext() {
		st, to := iter.Next()
		from := st.First
		symbol := charLabel(string(st.Second))
		edges[collection.NewPair(from, to)] = append(edges[collection.NewPair(from, to)], symbol)
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

func toStateRegexID(id RegexID) RegexID {
	if id == nonFinStateRegexID {
		return 0
	}

	return id
}
