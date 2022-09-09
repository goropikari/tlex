package automata_test

import (
	"testing"

	"github.com/goropikari/tlex/automata"
	"github.com/stretchr/testify/require"
)

func TestInterval_Overlap(t *testing.T) {
	tests := []struct {
		name     string
		x        automata.Interval
		y        automata.Interval
		expected bool
	}{
		{
			name:     "overlap",
			x:        automata.NewInterval(1, 3),
			y:        automata.NewInterval(2, 4),
			expected: true,
		},
		{
			name:     "overlap2",
			x:        automata.NewInterval(2, 4),
			y:        automata.NewInterval(1, 3),
			expected: true,
		},
		{
			name:     "overlap3",
			x:        automata.NewInterval(1, 3),
			y:        automata.NewInterval(3, 5),
			expected: true,
		},
		{
			name:     "non-overlap",
			x:        automata.NewInterval(1, 4),
			y:        automata.NewInterval(5, 6),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.x.Overlap(tt.y)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestInterval_Difference(t *testing.T) {
	tests := []struct {
		name     string
		x        automata.Interval
		y        automata.Interval
		expected []automata.Interval
	}{
		{
			name: "overlap",
			x:    automata.NewInterval(1, 3),
			y:    automata.NewInterval(2, 5),
			expected: []automata.Interval{
				automata.NewInterval(1, 1),
			},
		},
		{
			name: "overlap2",
			x:    automata.NewInterval(2, 5),
			y:    automata.NewInterval(1, 3),
			expected: []automata.Interval{
				automata.NewInterval(4, 5),
			},
		},
		{
			name: "overlap3",
			x:    automata.NewInterval(1, 7),
			y:    automata.NewInterval(3, 5),
			expected: []automata.Interval{
				automata.NewInterval(1, 2),
				automata.NewInterval(6, 7),
			},
		},
		{
			name:     "overlap4",
			x:        automata.NewInterval(3, 5),
			y:        automata.NewInterval(1, 7),
			expected: []automata.Interval{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.x.Difference(tt.y)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestInterval_Disjoin(t *testing.T) {
	tests := []struct {
		name     string
		given    []automata.Interval
		expected []automata.Interval
	}{
		{
			name: "different interval",
			given: []automata.Interval{
				automata.NewInterval(97, 99),
				automata.NewInterval(97, 100),
				automata.NewInterval(98, 108),
				automata.NewInterval(99, 99),
				automata.NewInterval(109, 200),
			},
			expected: []automata.Interval{
				automata.NewInterval(97, 97),
				automata.NewInterval(98, 98),
				automata.NewInterval(99, 99),
				automata.NewInterval(100, 100),
				automata.NewInterval(101, 108),
				automata.NewInterval(109, 200),
			},
		},
		{
			name: "same interval",
			given: []automata.Interval{
				automata.NewInterval(97, 99),
				automata.NewInterval(97, 99),
				automata.NewInterval(97, 99),
				automata.NewInterval(97, 99),
			},
			expected: []automata.Interval{
				automata.NewInterval(97, 99),
			},
		},
		{
			name: "points",
			given: []automata.Interval{
				automata.NewInterval(97, 97),
				automata.NewInterval(98, 98),
				automata.NewInterval(99, 99),
			},
			expected: []automata.Interval{
				automata.NewInterval(97, 97),
				automata.NewInterval(98, 98),
				automata.NewInterval(99, 99),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := automata.Disjoin(tt.given)

			require.Equal(t, tt.expected, got)

		})
	}

}
