package automata

// func (dfa DFA) Minimize() DFA {
// 	dfa = dfa.Totalize()
// 	dp := dfa.makeDP()

// 	stToID := make(map[State]int)
// 	sts := make([]State, 0)
// 	i := 0
// 	for v := range dfa.q {
// 		stToID[v] = i
// 		i++
// 		sts = append(sts, v)
// 	}

// 	uf := collection.NewUnionFind(len(dfa.q))
// 	for is := range dfa.q {
// 		for js := range dfa.q {
// 			if is.GetLabel() < js.GetLabel() {
// 				if dp[collection.NewTuple(is, js)] == -1 {
// 					uf.Unite(stToID[is], stToID[js])
// 				}
// 			}
// 		}
// 	}

// 	delta := make(DFATransition)
// 	for _, ru := range SupportedChars {
// 		for from := range dfa.q {
// 			tu := collection.NewTuple(sts[uf.Find(stToID[from])], ru)
// 			to := sts[uf.Find(stToID[dfa.delta[tu]])]
// 			delta[tu] = to
// 		}
// 	}

// 	q := collection.NewSet[State]()
// 	for _, st := range sts {
// 		ni := uf.Find(stToID[st])
// 		q.Insert(sts[ni])
// 	}

// 	finStates := collection.NewSet[State]()
// 	for st := range dfa.finStates {
// 		ni := uf.Find(stToID[st])
// 		finStates.Insert(sts[ni])
// 	}

// 	initState := sts[uf.Find(stToID[dfa.initState])]

// 	return NewDFA(q, delta, initState, finStates)
// }

// func (dfa DFA) makeDP() map[collection.Tuple[State, State]]int {
// 	initStateSet := collection.NewSet[State]()
// 	if !dfa.finStates.Contains(dfa.initState) {
// 		initStateSet.Insert(dfa.initState)
// 	}
// 	imStateSet := dfa.q.Difference(initStateSet).Difference(dfa.finStates)
// 	finStateSet := dfa.finStates

// 	dp := make(map[collection.Tuple[State, State]]int)
// 	for is := range dfa.q {
// 		for js := range dfa.q {
// 			if is.GetLabel() < js.GetLabel() {
// 				dp[collection.NewTuple(is, js)] = -1
// 			}
// 		}
// 	}
// 	for is := range initStateSet {
// 		for js := range imStateSet {
// 			if is.GetLabel() > js.GetLabel() {
// 				dp[collection.NewTuple(js, is)] = 0
// 			} else {
// 				dp[collection.NewTuple(is, js)] = 0
// 			}
// 		}
// 	}
// 	for is := range imStateSet {
// 		for js := range finStateSet {
// 			if is.GetLabel() > js.GetLabel() {
// 				dp[collection.NewTuple(js, is)] = 0
// 			} else {
// 				dp[collection.NewTuple(is, js)] = 0
// 			}

// 		}
// 	}
// 	for is := range initStateSet {
// 		for js := range finStateSet {
// 			if is.GetLabel() > js.GetLabel() {
// 				dp[collection.NewTuple(js, is)] = 0
// 			} else {
// 				dp[collection.NewTuple(is, js)] = 0
// 			}
// 		}
// 	}

// 	return dfa.fixPoint(dp)
// }

// // https://github.com/ganeshutah/Jove/blob/2f0e11794adc09bc8be08515917454f00c921c0e/jove/Def_DFA.py#L665
// func (dfa DFA) fixPoint(dp map[collection.Tuple[State, State]]int) map[collection.Tuple[State, State]]int {
// 	ok := true
// 	for ok {
// 		ok = false
// 		for pair := range dp {
// 			s0 := pair.First
// 			s1 := pair.Second
// 			for _, ru := range SupportedChars {
// 				ns0 := dfa.delta[collection.NewTuple(s0, ru)]
// 				ns1 := dfa.delta[collection.NewTuple(s1, ru)]
// 				if ns0 == ns1 {
// 					continue
// 				}
// 				if ns0.GetLabel() > ns1.GetLabel() {
// 					ns0, ns1 = ns1, ns0
// 				}
// 				p2 := collection.NewTuple(ns0, ns1)
// 				if v, ok := dp[p2]; ok {
// 					p1 := collection.NewTuple(s0, s1)
// 					if dp[p1] == -1 && v >= 0 {
// 						dp[p1] = dp[p2] + 1
// 						ok = true
// 						break
// 					}
// 				} else {
// 					panic(errors.New("invalid"))
// 				}
// 			}
// 		}
// 	}

// 	return dp
// }
