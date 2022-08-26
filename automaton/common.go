package automaton

const epsilon = 'ε'

// const SupportedChars = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ \t\n\r"

const SupportedChars = "ab"

type State struct {
	label string
}

func NewState(label string) State {
	return State{label: label}
}

func (st State) GetLabel() string {
	return st.label
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
