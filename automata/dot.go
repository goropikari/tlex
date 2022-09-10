package automata

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/golang/freetype/truetype"
	"github.com/goropikari/tlex/collection"
	"github.com/goropikari/tlex/utils/guid"
	"golang.org/x/image/font"
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
	nodes := make(map[StateID]*cgraph.Node)
	ii, si, fi := 0, 0, 0
	qiter := nfa.states.Iterator()
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
			e.SetLabel("ε")
			ii++
		}
		if nfa.finStates.Contains(s) {
			n.SetShape(cgraph.DoubleCircleShape)
			n.SetLabel(fmt.Sprintf("F%v_%v", fi, nfa.stIDToRegID.Get(s)))
			fi++
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v", si))
			si++
		}
		nodes[s] = n
	}

	edges := make(map[collection.Pair[StateID, StateID]]string)
	for from, mp := range nfa.trans.mp {
		for intv, tos := range mp {
			symbols := fmt.Sprintf("[%c-%c]", intv.L, intv.R)
			titer := tos.Iterator()
			for titer.HasNext() {
				to := titer.Next()
				if v, ok := edges[collection.NewPair(from, to)]; ok {
					edges[collection.NewPair(from, to)] = v + "\n" + symbols
				} else {
					edges[collection.NewPair(from, to)] = symbols
				}
			}
		}
	}
	for k, symbols := range edges {
		from := k.First
		to := k.Second
		e, err := graph.CreateEdge(symbols, nodes[from], nodes[to])
		if err != nil {
			panic(err)
		}
		e.SetLabel(charLabel(symbols))
	}

	for from, tos := range nfa.epsilonTrans.mp {
		iter := tos.Iterator()
		for iter.HasNext() {
			to := iter.Next()
			e, err := graph.CreateEdge("ε", nodes[from], nodes[to])
			if err != nil {
				panic(err)
			}
			e.SetLabel(charLabel("ε"))
		}
	}

	var buf bytes.Buffer
	// g.Render(graph, "dot", &buf)
	// s := buf.String()

	// if err := os.WriteFile("ex.dot", []byte(s), 0666); err != nil {
	// 	log.Fatal(err)
	// }
	// graph, err := graphviz.ParseBytes([]byte(s))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// g := graphviz.New()
	if err := g.RenderFilename(graph, graphviz.PNG, "ex.png"); err != nil {
		log.Fatal(err)
	}

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
	nodes := make(map[StateID]*cgraph.Node)
	ii, si, fi := 0, 0, 0
	for s := StateID(0); s < StateID(nfa.size); s++ {
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
			e.SetLabel("ε")
			ii++
		}
		if nfa.finStates.Contains(s) {
			n.SetShape(cgraph.DoubleCircleShape)
			n.SetLabel(fmt.Sprintf("F%v_%v", fi, nfa.stIDToRegID.Get(s)))
			fi++
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v", si))
			si++
		}
		nodes[s] = n
	}

	edges := make(map[collection.Pair[StateID, StateID]]string)
	for from, mp := range nfa.trans.mp {
		for intv, tos := range mp {
			symbols := fmt.Sprintf("[%c-%c]", intv.L, intv.R)
			titer := tos.iterator()
			for titer.HasNext() {
				to := titer.Next()
				if v, ok := edges[collection.NewPair(from, to)]; ok {
					edges[collection.NewPair(from, to)] = v + "\n" + symbols
				} else {
					edges[collection.NewPair(from, to)] = symbols
				}
			}
		}
	}
	for k, symbols := range edges {
		from := k.First
		to := k.Second
		e, err := graph.CreateEdge(symbols, nodes[from], nodes[to])
		if err != nil {
			panic(err)
		}
		e.SetLabel(charLabel(symbols))
	}

	for from, tos := range nfa.etrans.mp {
		iter := tos.iterator()
		for iter.HasNext() {
			to := iter.Next()
			e, err := graph.CreateEdge("ε", nodes[from], nodes[to])
			if err != nil {
				panic(err)
			}
			e.SetLabel(charLabel("ε"))
		}
	}

	var buf bytes.Buffer
	// g.Render(graph, "dot", &buf)
	// s := buf.String()

	// if err := os.WriteFile("ex.dot", []byte(s), 0666); err != nil {
	// 	log.Fatal(err)
	// }
	// graph, err := graphviz.ParseBytes([]byte(s))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// g := graphviz.New()
	if err := g.RenderFilename(graph, graphviz.PNG, "ex.png"); err != nil {
		log.Fatal(err)
	}

	return buf.String(), nil
}

func (dfa DFA) ToDot() (string, error) {
	g := graphviz.New()
	var ftBinary []byte
	if exists("/usr/share/fonts/opentype/ipaexfont-gothic/ipaexg.ttf") {
		ftBinary, _ = os.ReadFile("/usr/share/fonts/opentype/ipaexfont-gothic/ipaexg.ttf")
	} else if exists("/usr/share/fonts/OTF/ipaexm.ttf") {
		ftBinary, _ = os.ReadFile("/usr/share/fonts/OTF/ipaexm.ttf")
	} else {
		var err error
		ftBinary, err = os.ReadFile("./ipaexg00401/ipaexg.ttf")
		if err != nil {
			panic(err)
		}
	}
	ft, _ := truetype.Parse(ftBinary)
	g.SetFontFace(func(size float64) (font.Face, error) {
		opt := &truetype.Options{
			Size:              size,
			DPI:               0,
			Hinting:           0,
			GlyphCacheEntries: 0,
			SubPixelsX:        0,
			SubPixelsY:        0,
		}
		return truetype.NewFace(ft, opt), nil
	})

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
	nodes := make(map[StateID]*cgraph.Node)
	ii, si, fi := 0, 0, 0
	siter := dfa.states.Iterator()
	for siter.HasNext() {
		s := siter.Next()
		n, err := graph.CreateNode(fmt.Sprintf("%v", guid.New())) // assign unique node id
		if err != nil {
			return "", err
		}
		if dfa.initState == s {
			_, err := graph.CreateEdge(fmt.Sprintf("%v", guid.New()), start, n)
			if err != nil {
				return "", err
			}
			n.SetLabel(fmt.Sprintf("I%v", ii))
			ii++
		}
		if dfa.finStates.Contains(s) {
			n.SetShape(cgraph.DoubleCircleShape)
			rid := dfa.stIDToRegID.Get(s)
			n.SetLabel(fmt.Sprintf("F%v_%v", fi, rid))
			fi++
		} else {
			n.SetShape(cgraph.CircleShape)
			n.SetLabel(fmt.Sprintf("S%v", si))
			si++
		}
		nodes[s] = n
	}

	edges := make(map[collection.Pair[StateID, StateID]]string)
	for from, mp := range dfa.trans.delta {
		for intv, to := range mp {
			var lstr, rstr string
			lstr = fmt.Sprintf("%v", intv.L)
			rstr = fmt.Sprintf("%v", intv.R)
			symbols := fmt.Sprintf("[%s-%s]", lstr, rstr)
			p := collection.NewPair(from, to)
			if _, ok := edges[p]; ok {
				edges[p] = edges[p] + "\n" + symbols
			} else {
				edges[p] = symbols
			}
		}
	}
	for k, symbols := range edges {
		from := k.First
		to := k.Second
		e, err := graph.CreateEdge(symbols, nodes[from], nodes[to])
		if err != nil {
			panic(err)
		}
		e.SetLabel(charLabel(symbols))
	}

	var buf bytes.Buffer
	// g.Render(graph, "dot", &buf)
	// s := buf.String()

	// if err := os.WriteFile("ex.dot", []byte(s), 0666); err != nil {
	// 	log.Fatal(err)
	// }
	// graph, err := graphviz.ParseBytes([]byte(s))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// g := graphviz.New()
	if err := g.RenderFilename(graph, graphviz.PNG, "ex.png"); err != nil {
		log.Fatal(err)
	}

	return buf.String(), nil
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
