package automaton

import "github.com/goropikari/golex/collection"

const epsilon = 'ε'

type Transition map[collection.Tuple[State, rune]]collection.Set[State]

func (t Transition) Copy() Transition {
	delta := make(Transition)
	for k, v := range t {
		delta[k] = v.Copy()
	}

	return delta
}

type State struct {
	label string
}

func NewState(label string) State {
	return State{label: label}
}
