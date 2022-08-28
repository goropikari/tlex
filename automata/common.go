package automata

import (
	stdmath "math"
	"strings"

	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/math"
	"golang.org/x/exp/slices"
)

const epsilon = 'ε'

const SupportedChars = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ \t\n\r"

// const SupportedChars = "abcdefghijklmnopqrstuvwxyz+-*/.0123456789"

// const SupportedChars = "abc"

type TokenID int

type State struct {
	label   string
	tokenID TokenID
}

func NewState(label string) State {
	return State{label: label, tokenID: TokenID(stdmath.MaxInt)}
}

func (st State) GetLabel() string {
	return st.label
}

func (st State) GetTokenID() TokenID {
	if int(st.tokenID) == stdmath.MaxInt {
		return 0
	}
	return st.tokenID
}

func (st *State) SetTokenID(id TokenID) {
	st.tokenID = id
}

func NewStateSet(sts collection.Set[State]) State {
	label := labelConcat(sts)
	id := TokenID(stdmath.MaxInt)
	for st := range sts {
		id = math.Min(id, st.tokenID)
	}

	st := NewState(label)
	st.tokenID = id

	return st
}

func labelConcat(set collection.Set[State]) string {
	s := make([]string, 0, len(set))
	for v := range set {
		s = append(s, v.GetLabel())
	}
	slices.Sort(s)
	return strings.Join(s, "_")
}

func charLabel(s string) string {
	switch s {
	case "\t":
		return "Tab"
	case "\n":
		return "Newline"
	case "\r":
		return "CR"
	case " ":
		return "Space"
	case "\\":
		return "Backslash"
	default:
		return s
	}
}
