package automata

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	stdmath "math"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/math"
	"github.com/goropikari/golex/utils/guid"
)

type ImNFA struct {
	maxID       int
	stIDToRegID []TokenID
	delta       map[collection.Tuple[StateID, rune]]*StateSet
	initStates  *StateSet
	finStates   *StateSet
}

func (nfa ImNFA) ToDFA() DFA {
	numst := nfa.maxID + 1
	ecl := make([]*StateSet, numst)
	for i := 1; i <= nfa.maxID; i++ {
		b := nfa.eclosure(StateID(i))
		ecl[i] = b
	}

	initState := nfa.initStates.Copy()
	for i := 1; i <= nfa.maxID; i++ {
		sid := StateID(i)
		if nfa.initStates.Contains(sid) {
			initState = initState.Union(ecl[i])
		}
	}
	que := list.New() // list of *StateSet
	que.PushBack(initState)

	visited := map[Sha]*StateSet{}
	finStates := map[Sha]*StateSet{}
	initSha := initState.Sha256()
	if !initState.Intersection(nfa.finStates).IsEmpty() {
		finStates[initSha] = initState
	}
	visited[initSha] = initState

	delta := make(map[collection.Tuple[Sha, rune]]Sha)

	for que.Len() > 0 {
		top := que.Front()
		que.Remove(top)
		froms := top.Value.(*StateSet)

		for _, ru := range SupportedChars {
			tos := NewStateSet(numst)
			for fromID := 1; fromID <= nfa.maxID; fromID++ {
				fromStID := StateID(fromID)
				if !froms.Contains(fromStID) {
					continue
				}
				if nxs, ok := nfa.delta[collection.NewTuple(StateID(fromID), ru)]; ok {
					for nxID := 1; nxID <= nfa.maxID; nxID++ {
						nxStID := StateID(nxID)
						if nxs.Contains(nxStID) {
							tos = tos.Union(ecl[nxID])
						}
					}
				}
			}

			if tos.IsEmpty() {
				continue
			}
			to := tos.Sha256()
			if !tos.Intersection(nfa.finStates).IsEmpty() {
				finStates[to] = tos
			}
			delta[collection.NewTuple(froms.Sha256(), ru)] = to
			if _, ok := visited[to]; ok {
				continue
			}
			visited[to] = tos
			que.PushBack(tos)
		}
	}

	shaToStateID := map[Sha]StateID{}
	for key := range visited {
		shaToStateID[key] = StateID(guid.New())
	}
	stIDToState := map[StateID]State{}
	dfaStates := collection.NewSet[State]()
	for sha, id := range shaToStateID {
		st := NewState(id)
		if v, ok := visited[sha]; ok {
			if !v.Intersection(nfa.finStates).IsEmpty() {
				rid := TokenID(stdmath.MaxInt)
				for i := 1; i < numst; i++ {
					sid := StateID(i)
					if v.Contains(sid) {
						rid = math.Min(rid, nfa.stIDToRegID[StateID(i)])
					}
				}

				st.SetTokenID(TokenID(rid))
			}
		}
		dfaStates.Insert(st)
		stIDToState[id] = st
	}

	dfatrans := make(DFATransition)
	for pair, to := range delta {
		fromSha := pair.First
		ru := pair.Second
		dfatrans[collection.NewTuple(stIDToState[shaToStateID[fromSha]], ru)] = stIDToState[shaToStateID[to]]
	}

	dfaFinStates := collection.NewSet[State]()
	for s := range finStates {
		dfaFinStates.Insert(stIDToState[shaToStateID[s]])
	}

	return DFA{
		q:         dfaStates,
		delta:     dfatrans,
		initState: stIDToState[shaToStateID[initSha]],
		finStates: dfaFinStates,
	}
}

func (nfa ImNFA) eclosure(x StateID) *StateSet {
	que := list.New() // list of StateID
	que.PushBack(x)

	visited := NewStateSet(nfa.maxID + 1).Insert((x))
	closure := visited.Copy()
	for que.Len() > 0 {
		front := que.Front()
		que.Remove(front)
		top := front.Value.(StateID)

		if nxs, ok := nfa.delta[collection.NewTuple(top, epsilon)]; ok {
			closure = closure.Union(nxs)

			for nx := 1; nx <= nfa.maxID; nx++ {
				nxStID := StateID(nx)
				if nxs.Contains(nxStID) && !visited.Contains(nxStID) {
					visited = visited.Insert(nxStID)
					que.PushBack(nxStID)
				}
			}
		}
	}

	return closure
}

func (nfa ImNFA) ToDot() (string, error) {
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
		st.SetTokenID(nfa.stIDToRegID[StateID(id)])
		nodes[st] = n
	}

	for st, qs := range nfa.delta {
		from := st.First
		symbol := string(st.Second)
		fromst := NewState(from)
		fromst.SetTokenID(nfa.stIDToRegID[from])
		for id := 1; id <= nfa.maxID; id++ {
			sid := StateID(id)
			if !qs.Contains(sid) {
				continue
			}
			tost := NewState(sid)
			tost.SetTokenID(nfa.stIDToRegID[id])
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
