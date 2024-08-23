package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Head(t *testing.T) {
	{
		h := NewHeap[int]()
		h.Put(0)
		h.Put(1)
		h.Put(2)

		s := h.List()

		n := len(s)
		var exp = make([]int, n)
		exp[n-1], exp[n-2] = 2, 1
		require.Equal(t, exp, s)
	}

	{
		h := NewHeap[int]()

		n := len(h.List())
		for i := 0; i < n+1; i++ {
			h.Put(i)
		}

		s := h.List()

		var exp = make([]int, 0, n)
		for i := 0; i < n; i++ {
			exp = append(exp, i+1)
		}
		require.Equal(t, exp, s)
	}
}
