package lexer

// dfa から生成される Lexer の想像図

import (
	"errors"
	"fmt"
)

// embedded code がここに入る
const (
	Keyword int = iota + 1
	Digit
	Identifier
)

// embedded code 終わり

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
	0, // state0 は BH state
	0, // 1
	1, // 2
	0, // 3
	3, // 4
	3, // 5
	2, // 6
}

// 生成する
var yyFinStates = map[yyStateID]struct{}{
	2: {},
	4: {},
	5: {},
	6: {},
}

// 生成する
var yyTransitionTable = map[yyStateID]map[rune]yyStateID{
	1: {
		'a': 2,
		'b': 4,
	},
	2: {
		'a': 3,
		'b': 5,
	},
	3: {
		'a': 3,
		'b': 4,
	},
	4: {
		'b': 4,
	},
	5: {
		'b': 6,
	},
	6: {
		'b': 4,
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
		currStateID: 1, // ここは後で埋め込む or init state id を 1 になるようにする
	}
}

func (yylex *yyLexer) currRune() rune {
	if yylex.currPos >= yylex.length {
		return 0
	}

	return yylex.data[yylex.currPos]
}

func (yylex *yyLexer) Next() (int, error) {
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
			case 1: // case 1 以降は生成する
				{
					// 埋め込み
					fmt.Println("state: 1")
					return Keyword, nil
				}
			case 2:
				{
					// 埋め込み
					fmt.Println("state: 2")
					return Digit, nil
				}
			case 3:
				{
					// 埋め込み
					fmt.Println("state: 3")
					return Identifier, nil
				}
				// ... regex pattern の数だけ生成
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
