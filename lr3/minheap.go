package main

import (
	"fmt"
	"math"
	"strings"
)

type Comparer interface {
	Compare(Comparer) int
}

// Коллекция элементов с методом сравнения.
// Может использоваться в min и max MinHeap
type MinHeap []Comparer

// min-heap
func NewMinHeap() MinHeap {
	heap := make(MinHeap, 0)

	return heap
}

func (hptr *MinHeap) Push(i Comparer) {
	h := *hptr

	h = append(h, i)

	idx := len(h) - 1

	for idx > 0 {
		parentIdx := idx / 2

		if h[idx].Compare(h[parentIdx]) > 0 {
			break
		}

		h[idx], h[parentIdx] = h[parentIdx], h[idx]
		idx = parentIdx
	}

	*hptr = h
}

func (hptr *MinHeap) Pop() Comparer {
	h := *hptr
	n := len(h)

	// пустая куча
	if n == 0 {
		return nil
	}

	idx := 0
	root := h[idx]

	v := h[n-1]

	for {
		childIdx := idx*2 + 1
		if childIdx >= n {
			break // больше нет детей, выход
		}

		if childIdx+1 < n && h[childIdx].Compare(h[childIdx+1]) > 0 {
			childIdx += 1
		}

		if v.Compare(h[childIdx]) < 0 {
			break
		}

		h[idx] = h[childIdx]
		idx = childIdx
	}

	h[idx] = v

	*hptr = h[:len(h)-1]

	return root
}

func (hptr *MinHeap) PrintHeap() {
	heapy := MinHeap(*hptr)

	size := len(heapy)

	maxDepth := int(math.Log(float64(size)) / math.Log(2))

	var hs strings.Builder

	for d := 0; d <= maxDepth; d++ {
		layerLength := int(math.Pow(2, float64(d)))

		var line strings.Builder

		for i := layerLength; i < int(math.Pow(2, float64(d+1))); i++ {
			if d != maxDepth {
				line.WriteString(strings.Repeat("  ", int(math.Pow(2, float64(maxDepth-d)))))
			}

			loops := maxDepth - d

			if loops >= 2 {
				loops -= 2
				for loops >= 0 {
					line.WriteString(strings.Repeat("  ", int(math.Pow(2, float64(loops)))))
					loops--
				}
			}

			if i <= size {
				line.WriteString(fmt.Sprintf("%-4d", heapy[i-1]))
			} else {
				line.WriteString(" -- ")
			}

			line.WriteString(strings.Repeat("  ", int(math.Pow(2, float64(maxDepth-d)))))

			loops = maxDepth - d
			if loops >= 2 {
				loops -= 2
				for loops >= 0 {
					line.WriteString(strings.Repeat("  ", int(math.Pow(2, float64(loops)))))
					loops--
				}
			}
		}

		hs.WriteString(line.String() + "\n")
	}

	fmt.Println(hs.String())
}

func (hptr *MinHeap) Traverse(idx int) {
	h := MinHeap(*hptr)

	if idx >= len(h) {
		return
	}

	fmt.Printf("%d ->", h[idx])
	hptr.Traverse((2 * idx) + 1)
	hptr.Traverse((2 * idx) + 2)
}

type Int int

func (a Int) Compare(b Comparer) int {
	bint := b.(Int)

	if a < bint {
		return -1
	} else if a == bint {
		return 0
	} else {
		return 1
	}
}
