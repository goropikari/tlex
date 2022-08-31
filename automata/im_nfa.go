package automata

import (
	"bytes"
	"container/list"
	"crypto/sha256"
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
	delta       map[collection.Tuple[StateID, rune]]collection.Bitset
	initStates  collection.Bitset
	finStates   collection.Bitset
}

func (nfa ImNFA) ToDFA() DFA {
	numst := nfa.maxID + 1
	ecl := make([]collection.Bitset, numst)
	for i := 1; i <= nfa.maxID; i++ {
		b := nfa.eclosure(StateID(i))
		ecl[i] = b
	}

	initStateBitset := nfa.initStates.Copy()
	for i := 1; i <= nfa.maxID; i++ {
		if nfa.initStates.Contains(i) {
			initStateBitset = initStateBitset.Union(ecl[i])
		}
	}
	que := list.New() // list of collection.Bitset
	que.PushBack(initStateBitset)

	memo := map[Sha]collection.Bitset{}
	dfaFinStateSet := map[Sha]collection.Bitset{}
	initSha := buildSha256(initStateBitset)
	if !initStateBitset.Intersection(nfa.finStates).IsZero() {
		dfaFinStateSet[initSha] = initStateBitset
	}
	memo[initSha] = initStateBitset

	dfaDelta := make(map[collection.Tuple[Sha, rune]]Sha)

	for que.Len() > 0 {
		top := que.Front()
		que.Remove(top)
		froms := top.Value.(collection.Bitset)

		for _, ru := range SupportedChars {
			tos := collection.NewBitset(numst)
			for fromStID := 1; fromStID <= nfa.maxID; fromStID++ {
				if !froms.Contains(fromStID) {
					continue
				}
				if nx, ok := nfa.delta[collection.NewTuple(StateID(fromStID), ru)]; ok {
					for nxStID := 1; nxStID <= nfa.maxID; nxStID++ {
						if nx.Contains(nxStID) {
							tos = tos.Union(ecl[nxStID])
						}
					}
				}
			}

			if tos.IsZero() {
				continue
			}
			to := buildSha256(tos)
			if !tos.Intersection(nfa.finStates).IsZero() {
				dfaFinStateSet[to] = tos
			}
			dfaDelta[collection.NewTuple(buildSha256(froms), ru)] = to
			if _, ok := memo[to]; ok {
				continue
			}
			memo[to] = tos
			que.PushBack(tos)
		}
	}

	shaToStateID := map[Sha]StateID{}
	for key := range memo {
		shaToStateID[key] = StateID(guid.New())
	}
	stIDToState := map[StateID]State{}
	dfaStates := collection.NewSet[State]()
	for sha, id := range shaToStateID {
		st := NewState(id)
		if v, ok := memo[sha]; ok {
			if !v.Intersection(nfa.finStates).IsZero() {
				rid := TokenID(stdmath.MaxInt)
				for i := 1; i < numst; i++ {
					if v.Contains(i) {
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
	for pair, to := range dfaDelta {
		fromSha := pair.First
		ru := pair.Second
		dfatrans[collection.NewTuple(stIDToState[shaToStateID[fromSha]], ru)] = stIDToState[shaToStateID[to]]
	}

	dfaFinStates := collection.NewSet[State]()
	for s := range dfaFinStateSet {
		dfaFinStates.Insert(stIDToState[shaToStateID[s]])
	}

	return DFA{
		q:         dfaStates,
		delta:     dfatrans,
		initState: stIDToState[shaToStateID[initSha]],
		finStates: dfaFinStates,
	}
}

func (nfa ImNFA) eclosure(x StateID) collection.Bitset {
	numst := nfa.maxID + 1
	que := list.New() // list of Sha of StateID
	// shaID := sha256.Sum256(collection.NewBitset(numst).Up(int(x)).x)
	que.PushBack(x)

	visited := collection.NewBitset(numst).Up(int(x))
	closure := visited.Copy()
	for que.Len() > 0 {
		front := que.Front()
		que.Remove(front)
		top := front.Value.(StateID)

		if nxs, ok := nfa.delta[collection.NewTuple(top, epsilon)]; ok {
			closure = closure.Union(nxs)

			for nx := 0; nx < numst; nx++ {
				if nxs.Contains(nx) && !visited.Contains(nx) {
					visited = visited.Up(nx)
					// shaID := sha256.Sum256(collection.NewBitset(numst).Up(int(nx)).x)
					que.PushBack(StateID(nx))
				}
			}
		}
	}

	return closure
}

func bitset(sz int, x int) (Sha, collection.Bitset) {
	b := collection.NewBitset(sz + 1).Up(x)
	s := sha256.Sum256(b.Bytes())
	return s, b
}

func buildSha256(bs collection.Bitset) Sha {
	return sha256.Sum256(bs.Bytes())
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
	for stid := 1; stid <= nfa.maxID; stid++ {
		n, err := graph.CreateNode(fmt.Sprintf("%v", guid.New())) // assign unique node id
		if err != nil {
			return "", err
		}
		if nfa.initStates.Contains(stid) {
			e, err := graph.CreateEdge(fmt.Sprintf("%v", guid.New()), start, n)
			if err != nil {
				return "", err
			}
			n.SetLabel(fmt.Sprintf("I%v", ii))
			e.SetLabel(string(epsilon))
			ii++
		}
		if nfa.finStates.Contains(stid) {
			n.SetShape(cgraph.DoubleCircleShape)
			n.SetLabel(fmt.Sprintf("F%v", fi))
			fi++
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v", si))
			si++
		}
		st := NewState(StateID(stid))
		st.SetTokenID(nfa.stIDToRegID[StateID(stid)])
		nodes[st] = n
	}

	for st, qs := range nfa.delta {
		from := st.First
		symbol := string(st.Second)
		fromst := NewState(from)
		fromst.SetTokenID(nfa.stIDToRegID[from])
		for to := 1; to <= nfa.maxID; to++ {
			if !qs.Contains(to) {
				continue
			}
			tost := NewState(StateID(to))
			tost.SetTokenID(nfa.stIDToRegID[to])
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
