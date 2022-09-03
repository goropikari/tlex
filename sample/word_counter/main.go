package main

import (
	"errors"
	"fmt"
)

var nc = 0
var nw = 0
var nl = 0

type yyStateID = int
type yyRegexID = int

var YYText string

var (
	ErrYYScan = errors.New("failed to scan")
	EOF       = errors.New("EOF")
)

// state id to regex id
var yyStateIDToRegexID = []yyRegexID{
	0, // state 0 は BH state
	1,
	1,
	2,
	3,
}

var yyFinStates = map[yyStateID]struct{}{
	1: {},
	2: {},
	3: {},
	4: {},
}

var yyTransitionTable = map[yyStateID]map[rune]yyStateID{
	1: {
		40:  2,
		41:  2,
		48:  2,
		49:  2,
		50:  2,
		51:  2,
		52:  2,
		53:  2,
		54:  2,
		55:  2,
		56:  2,
		57:  2,
		65:  2,
		66:  2,
		67:  2,
		68:  2,
		69:  2,
		70:  2,
		71:  2,
		72:  2,
		73:  2,
		74:  2,
		75:  2,
		76:  2,
		77:  2,
		78:  2,
		79:  2,
		80:  2,
		81:  2,
		82:  2,
		83:  2,
		84:  2,
		85:  2,
		86:  2,
		87:  2,
		88:  2,
		89:  2,
		90:  2,
		97:  2,
		98:  2,
		99:  2,
		100: 2,
		101: 2,
		102: 2,
		103: 2,
		104: 2,
		105: 2,
		106: 2,
		107: 2,
		108: 2,
		109: 2,
		110: 2,
		111: 2,
		112: 2,
		113: 2,
		114: 2,
		115: 2,
		116: 2,
		117: 2,
		118: 2,
		119: 2,
		120: 2,
		121: 2,
		122: 2,
		32:  3,
		9:   3,
		10:  4,
		13:  3,
	},
	2: {
		40:  2,
		41:  2,
		48:  2,
		49:  2,
		50:  2,
		51:  2,
		52:  2,
		53:  2,
		54:  2,
		55:  2,
		56:  2,
		57:  2,
		65:  2,
		66:  2,
		67:  2,
		68:  2,
		69:  2,
		70:  2,
		71:  2,
		72:  2,
		73:  2,
		74:  2,
		75:  2,
		76:  2,
		77:  2,
		78:  2,
		79:  2,
		80:  2,
		81:  2,
		82:  2,
		83:  2,
		84:  2,
		85:  2,
		86:  2,
		87:  2,
		88:  2,
		89:  2,
		90:  2,
		97:  2,
		98:  2,
		99:  2,
		100: 2,
		101: 2,
		102: 2,
		103: 2,
		104: 2,
		105: 2,
		106: 2,
		107: 2,
		108: 2,
		109: 2,
		110: 2,
		111: 2,
		112: 2,
		113: 2,
		114: 2,
		115: 2,
		116: 2,
		117: 2,
		118: 2,
		119: 2,
		120: 2,
		121: 2,
		122: 2,
	},
	3: {
		32: 3,
		9:  3,
		13: 3,
	},
}

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
	YYText      string
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
		currStateID: 1, // init state id is 1.
	}
}

func (yylex *yyLexer) currRune() rune {
	if yylex.currPos >= yylex.length {
		return 0
	}

	return yylex.data[yylex.currPos]
}

func (yylex *yyLexer) Next() (int, error) {
yystart:
	if yylex.currPos >= yylex.length {
		return 0, EOF
	}
	for yylex.currPos <= yylex.length {
		yyNxStID := yyNextStep(yylex.currStateID, yylex.currRune())
		if yyNxStID == 0 {
			yylex.YYText = string(yylex.data[yylex.beginPos : yylex.finPos+1])
			YYText = yylex.YYText
			yyNewCurrPos := yylex.finPos + 1
			yylex.beginPos = yyNewCurrPos
			yylex.finPos = yyNewCurrPos
			yylex.currPos = yyNewCurrPos
			yylex.currStateID = 1

			regexID := yylex.finRegexID
			yylex.finRegexID = 0
			switch regexID {
			case 0:
				return 0, ErrYYScan
			case 1:
				{
					nc += len(YYText)
					nw++
				}
				goto yystart
			case 2:
				{
					nc++
				}
				goto yystart
			case 3:
				{
					nl++
					nc++
				}
				goto yystart

			default:
				return 0, ErrYYScan
			}
		}
		if _, ok := yyFinStates[yyNxStID]; ok {
			yylex.finPos = yylex.currPos
			yylex.finRegexID = yyStateIDToRegexID[yyNxStID]
		}
		yylex.currStateID = yyNxStID
		yylex.currPos++
	}

	return 0, ErrYYScan
}

// This part is optional
func main() {
	program := `hello world
hello tlex
`
	fmt.Println(program)
	fmt.Println("-----------------")

	lex := New(program)
	for {
		_, err := lex.Next()
		if err != nil {
			break
		}
	}
	fmt.Printf("number of lines: %d\n", nl)
	fmt.Printf("number of words: %d\n", nw)
	fmt.Printf("number of chars: %d\n", nc)
}
