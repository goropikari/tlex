package automata

import (
	"crypto/sha256"
	stdmath "math"
)

const epsilon = 'ε'

const SupportedChars = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ \t\n\r"

// const SupportedChars = "abcdefghijklmnopqrstuvwxyz+-*/.0123456789"

// const SupportedChars = "abc"

type TokenID int
type StateID int
type Sha = [sha256.Size]byte

type State struct {
	id      StateID
	tokenID TokenID
}

func NewState(id StateID) State {
	return State{id: id, tokenID: TokenID(stdmath.MaxInt)}
}

func (st State) GetID() StateID {
	return st.id
}

func (st State) GetTokenID() TokenID {
	if int(st.tokenID) == stdmath.MaxInt {
		return 0
	}
	return st.tokenID
}

func (st State) GetRawTokenID() TokenID {
	return st.tokenID
}

func (st *State) SetTokenID(id TokenID) {
	st.tokenID = id
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
