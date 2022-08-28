package main

import (
	"errors"
	"fmt"
	"log"
)

type Type = int

const (
	State1 Type = iota + 1
	State2
	State3
	Other
)

type yyStateID = int
type yyRegexID = int

var YYtext string
var (
	ErrYYScan = errors.New("failed to scan")
	EOF       = errors.New("EOF")
)

// 生成する
// state id to regex id
var yyStateIDToRegexID = []yyRegexID{
	0, // state 0 は BH state
	0,
	2,
	1,
	0,
	3,
	3,
}

// 生成する
var yyFinStates = map[yyStateID]struct{}{
	2: {},
	3: {},
	5: {},
	6: {},
}

// 生成する
var yyTransitionTable = map[yyStateID]map[rune]yyStateID{
	1: {
		97: 3,
		98: 6,
	},
	2: {
		98: 6,
	},
	3: {
		97: 4,
		98: 5,
	},
	4: {
		97: 4,
		98: 6,
	},
	5: {
		98: 2,
	},
	6: {
		98: 6,
	},
}

// ここは固定値
func yyNextStep(id yyStateID, ru rune) yyStateID {
	if mp, ok := yyTransitionTable[id]; ok {
		return mp[ru]
	}

	return 0
}

type yyLexer struct {
	data        []rune
	length      int
	beginPos    int
	finPos      int
	currPos     int
	finRegexID  int
	currStateID yyStateID
}

func New(data string) *yyLexer {
	runes := []rune(data)
	return &yyLexer{
		data:        runes,
		length:      len(runes),
		beginPos:    0,
		finPos:      0,
		currPos:     0,
		finRegexID:  0,
		currStateID: 1, // init state id を 1 になるようにする
	}
}

func (yylex *yyLexer) currRune() rune {
	if yylex.currPos >= yylex.length {
		return 0
	}

	return yylex.data[yylex.currPos]
}

func (yylex *yyLexer) Next() (int, error) {
start:
	if yylex.currPos >= yylex.length {
		return 0, EOF
	}
	for yylex.currPos <= yylex.length {
		nxStID := yyNextStep(yylex.currStateID, yylex.currRune())
		if nxStID == 0 {
			YYtext = string(yylex.data[yylex.beginPos : yylex.finPos+1])
			newCurrPos := yylex.finPos + 1
			yylex.beginPos = newCurrPos
			yylex.finPos = newCurrPos
			yylex.currPos = newCurrPos
			yylex.currStateID = 1

			regexID := yylex.finRegexID
			yylex.finRegexID = 0
			switch regexID {
			case 0:
				return 0, ErrYYScan
			case 1:
				{
					return State1, nil
				}
				goto start
			case 2:
				{
					return State2, nil
				}
				goto start
			case 3:
				{
					return State3, nil
				}
				goto start
			case 4:
				{
					return Other, nil
				}
				goto start

			default:
				return 0, ErrYYScan
			}
		}
		if _, ok := yyFinStates[nxStID]; ok {
			yylex.finPos = yylex.currPos
			yylex.finRegexID = yyStateIDToRegexID[nxStID]
		}
		yylex.currStateID = nxStID
		yylex.currPos++
	}

	return 0, ErrYYScan
}

func main() {
	lex := New("ababba")
	for {
		n, err := lex.Next()
		if err != nil {
			log.Fatal(err)
			return
		}
		switch n {
		case State1:
			fmt.Println(State1, YYtext)
		case State2:
			fmt.Println(State2, YYtext)
		case State3:
			fmt.Println(State3, YYtext)
		default:
			fmt.Println(n, YYtext)
		}
	}
}
