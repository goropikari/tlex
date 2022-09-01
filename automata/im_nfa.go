package automata

import (
	"bytes"
	"container/list"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/math"
	"github.com/goropikari/golex/utils/guid"
)

type ImdNFATransition map[collection.Tuple[StateID, rune]]*StateSet

func (trans ImdNFATransition) step(x StateID, ru rune) (*StateSet, bool) {
	nxs, ok := trans[collection.NewTuple(x, ru)]
	return nxs, ok
}

type ImdNFA struct {
	maxID       int
	stIDToRegID []TokenID
	delta       ImdNFATransition
	initStates  *StateSet
	finStates   *StateSet
}

func (nfa ImdNFA) buildEClosures() []*StateSet {
	ecl := make([]*StateSet, nfa.numst())
	iter := nfa.iterator()
	for iter.HasNext() {
		sid := iter.Next()
		b := nfa.eclosure(sid)
		ecl[sid] = b
	}

	return ecl
}

func (nfa ImdNFA) numst() int {
	// +1 means black hole
	return nfa.maxID + 1
}

func (nfa ImdNFA) genStateID() StateID {
	return StateID(guid.New())
}

func (nfa ImdNFA) calRegID(ss *StateSet) TokenID {
	regID := nonFinStateTokenID
	iter := ss.iterator()
	for iter.HasNext() {
		sid := iter.Next()
		regID = math.Min(regID, nfa.stIDToRegID[sid])
	}

	return regID
}

func (nfa ImdNFA) ToDFA() DFA {
	states, delta, initState, finStates := nfa.subsetConstruction()

	dfaStates := collection.NewSet[State]()
	ssToSt := NewStateSetDict[State]()
	siter := states.iterator()
	for siter.HasNext() {
		ss, newSid := siter.Next()
		regID := nfa.calRegID(ss)
		st := NewStateWithTokenID(newSid, regID)
		dfaStates.Insert(st)
		ssToSt.Set(ss, st)
	}

	diter := delta.iterator()
	dfaDelta := make(DFATransition)
	for diter.HasNext() {
		fromSs, mp := diter.Next()
		for ru, toSs := range mp {
			fromSt, _ := ssToSt.Get(fromSs)
			toSt, _ := ssToSt.Get(toSs)
			dfaDelta[collection.NewTuple(fromSt, ru)] = toSt
		}
	}

	dfaInitState, _ := ssToSt.Get(initState)

	dfaFinStates := collection.NewSet[State]()
	fiter := finStates.iterator()
	for fiter.HasNext() {
		ss, _ := fiter.Next()
		st, _ := ssToSt.Get(ss)
		dfaFinStates.Insert(st)
	}

	return NewDFA(dfaStates, dfaDelta, dfaInitState, dfaFinStates)
}

func (nfa ImdNFA) subsetConstruction() (states *StateSetDict[StateID], delta *StateSetDict[map[rune]*StateSet], initState *StateSet, finStates *StateSetDict[Nothing]) {
	ecl := nfa.buildEClosures()

	initState = nfa.initStates.Copy()
	initIter := initState.iterator()
	for initIter.HasNext() {
		sid := initIter.Next()
		initState = initState.Union(ecl[sid])
	}

	visited := NewStateSetDict[StateID]()
	finStates = NewStateSetDict[Nothing]()
	if initState.Intersection(nfa.finStates).IsAny() {
		finStates.Set(initState, nothing)
	}
	visited.Set(initState, nfa.genStateID())

	delta = NewStateSetDict[map[rune]*StateSet]()

	que := list.New() // list of *StateSet
	que.PushBack(initState)
	for que.Len() > 0 {
		top := que.Front()
		que.Remove(top)
		from := top.Value.(*StateSet)

		for _, ru := range SupportedChars {
			tos := NewStateSet(nfa.numst())
			fromIter := from.iterator()
			for fromIter.HasNext() {
				fromStID := fromIter.Next()
				if nxs, ok := nfa.delta.step(fromStID, ru); ok {
					nxsIter := nxs.iterator()
					for nxsIter.HasNext() {
						nxStID := nxsIter.Next()
						tos = tos.Union(ecl[nxStID])
					}
				}
			}

			if tos.IsEmpty() {
				continue
			}
			if tos.Intersection(nfa.finStates).IsAny() {
				finStates.Set(tos, nothing)
			}
			if v, ok := delta.Get(from); ok {
				v[ru] = tos
				delta.Set(from, v)
			} else {
				mp := map[rune]*StateSet{}
				mp[ru] = tos
				delta.Set(from, mp)
			}
			if visited.Contains(tos) {
				continue
			}
			visited.Set(tos, nfa.genStateID())
			que.PushBack(tos)
		}
	}

	return visited, delta, initState, finStates
}

func (nfa ImdNFA) eclosure(x StateID) *StateSet {
	que := list.New() // list of StateID
	que.PushBack(x)

	visited := NewStateSet(nfa.maxID + 1).Insert((x))
	closure := visited.Copy()
	for que.Len() > 0 {
		front := que.Front()
		que.Remove(front)
		top := front.Value.(StateID)

		if nxs, ok := nfa.delta.step(top, epsilon); ok {
			closure = closure.Union(nxs)
			nxsIter := nxs.iterator()
			for nxsIter.HasNext() {
				nxStID := nxsIter.Next()
				if !visited.Contains(nxStID) {
					visited = visited.Insert(nxStID)
					que.PushBack(nxStID)
				}
			}
		}
	}

	return closure
}

func (nfa ImdNFA) iterator() *allStateIDIterator {
	return newAllStateIDIterator(nfa.maxID)
}

type allStateIDIterator struct {
	maxID  int
	currID int
}

func newAllStateIDIterator(maxID int) *allStateIDIterator {
	return &allStateIDIterator{
		maxID:  maxID,
		currID: 1, // StateID = 0 is blackhole state
	}
}

func (iter *allStateIDIterator) HasNext() bool {
	return iter.currID <= iter.maxID
}

func (iter *allStateIDIterator) Next() StateID {
	ret := StateID(iter.currID)
	iter.currID++
	return ret
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
