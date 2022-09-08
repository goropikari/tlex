package collection_test

import (
	"testing"

	"github.com/goropikari/tlex/collection"
	"github.com/stretchr/testify/require"
)

func TestPriorityQueue(t *testing.T) {
	tests := []struct {
		name     string
		given    []int
		less     func(x, y int) bool
		expected []int
	}{
		{
			name:     "descending",
			given:    []int{2, 5, 1, 9, 3, 9},
			less:     func(x, y int) bool { return x < y },
			expected: []int{9, 9, 5, 3, 2, 1},
		},
		{
			name:     "ascending",
			given:    []int{2, 5, 1, 9, 3, 9},
			less:     func(x, y int) bool { return x > y },
			expected: []int{1, 2, 3, 5, 9, 9},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			pq := collection.NewPriorityQueue(tt.less)
			for _, v := range tt.given {
				pq.Push(v)
			}

			got := make([]int, 0, len(tt.expected))
			for !pq.IsEmpty() {
				x := pq.Top()
				pq.Pop()
				got = append(got, x)
			}

			require.Equal(t, tt.expected, got)
		})
	}
}
