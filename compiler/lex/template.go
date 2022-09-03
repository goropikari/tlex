package lex

const tmpl = `
package {{ .PackageName }}

import (
	"errors"
	"fmt"
)

{{ .EmbeddedTmpl }}

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
    {{ .StateIDToRegexIDTmpl }}
}

var yyFinStates = map[yyStateID]struct{}{
    {{ .FinStatesTmpl }}
}

var yyTransitionTable = map[yyStateID]map[rune]yyStateID{
    {{ .TransitionTableTmpl }}
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
            {{ .RegexActionsTmpl }}
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

{{ .UserCodeTmpl }}
`
