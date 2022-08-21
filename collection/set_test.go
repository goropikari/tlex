package collection_test

import (
	"testing"

	"github.com/goropikari/golex/collection"
	"github.com/stretchr/testify/require"
)

func TestSet_Insert(t *testing.T) {
	tests := []struct {
		name     string
		given    []int
		expected collection.Set[int]
	}{
		{
			name:  "Insert",
			given: []int{1, 2, 3},
			expected: collection.Set[int]{
				1: struct{}{},
				2: struct{}{},
				3: struct{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := collection.NewSet[int]()
			for _, v := range tt.given {
				s.Insert(v)
			}

			require.Equal(t, tt.expected, s)
		})
	}
}

func TestSet_Erase(t *testing.T) {
	tests := []struct {
		name     string
		given    collection.Set[int]
		erase    []int
		expected collection.Set[int]
	}{
		{
			name: "Erase",
			given: collection.Set[int]{
				1: struct{}{},
				2: struct{}{},
				3: struct{}{},
			},
			erase: []int{1, 2},
			expected: collection.Set[int]{
				3: struct{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.erase {
				tt.given.Erase(v)
			}

			require.Equal(t, tt.expected, tt.given)
		})
	}
}

func TestSet_Difference(t *testing.T) {
	tests := []struct {
		name     string
		lhs      collection.Set[int]
		rhs      collection.Set[int]
		expected collection.Set[int]
	}{
		{
			name: "Difference",
			lhs: collection.Set[int]{
				1: struct{}{},
				2: struct{}{},
				3: struct{}{},
			},
			rhs: collection.Set[int]{
				1: struct{}{},
				2: struct{}{},
			},
			expected: collection.Set[int]{
				3: struct{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.lhs.Difference(tt.rhs))
		})
	}
}

func TestSet_Intersection(t *testing.T) {
	tests := []struct {
		name     string
		lhs      collection.Set[int]
		rhs      collection.Set[int]
		expected collection.Set[int]
	}{
		{
			name: "Intersection",
			lhs: collection.Set[int]{
				1: struct{}{},
				2: struct{}{},
				3: struct{}{},
			},
			rhs: collection.Set[int]{
				1: struct{}{},
				2: struct{}{},
			},
			expected: collection.Set[int]{
				1: struct{}{},
				2: struct{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.lhs.Intersection(tt.rhs))
		})
	}
}

func TestSet_Union(t *testing.T) {
	tests := []struct {
		name     string
		lhs      collection.Set[int]
		rhs      collection.Set[int]
		expected collection.Set[int]
	}{
		{
			name: "Union",
			lhs: collection.Set[int]{
				2: struct{}{},
				3: struct{}{},
			},
			rhs: collection.Set[int]{
				1: struct{}{},
				2: struct{}{},
			},
			expected: collection.Set[int]{
				1: struct{}{},
				2: struct{}{},
				3: struct{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.lhs.Union(tt.rhs))
		})
	}
}
