package main

import (
	"fmt"
	"math"
	"strings"
)

// Коллекция элементов с методом сравнения.
// Может использоваться в min и max MinHeap
type MinHeap struct {
	heap []int
}

// min-heap
func NewMinHeap() MinHeap {
	heap := make([]int, 0)

	return MinHeap{
		heap: heap,
	}
}

func (mh *MinHeap) Push(newElem int) {
	h := mh.heap

	// Добавляем новый элемент в конец массива
	// и сохраняем его индекс
	h = append(h, newElem)
	idx := len(h) - 1

	// пока не прошли все родительские элементы
	for idx > 0 {
		// родительский элемент
		parentIdx := idx / 2

		// новый элемент больше родительского
		// оставляем его на месте и выходим
		if h[idx] > h[parentIdx] {
			break
		}

		// меняем местами
		h[idx], h[parentIdx] = h[parentIdx], h[idx]
		idx = parentIdx
	}

	mh.heap = h
}

func (mh *MinHeap) Pop() int {
	h := mh.heap

	// размер кучи
	n := len(h)

	// пустая куча
	if n == 0 {
		return -1
	}

	// сохраняем корневой элемент
	idx := 0
	root := h[idx]

	// последний элемент в куче
	v := h[n-1]

	for {
		// индекс левого потомка текущего узла
		childIdx := idx*2 + 1

		// если левый потомок находится за пределами кучи
		if childIdx >= n {
			break
		}

		// если есть правый потомок и он меньше левого, выбираем его
		if childIdx+1 < n && h[childIdx] > h[childIdx+1] {
			childIdx += 1
		}

		// если значение последнего элемента меньше чем значение выбранного потомка, выходим
		if v < h[childIdx] {
			break
		}

		// перемещаем значение потомка в текущий узел
		// и обновляем индекс текущего узла
		h[idx] = h[childIdx]
		idx = childIdx
	}

	h[idx] = v

	// обновляем кучу, удаляя последний элемент
	mh.heap = h[:len(h)-1]

	// возвращаем удаленное минимальное значение
	return root
}

func (mh *MinHeap) Traverse(idx int) {
	h := mh.heap

	if idx >= len(h) {
		return
	}

	fmt.Printf("%d ->", h[idx])
	mh.Traverse((2 * idx) + 1)
	mh.Traverse((2 * idx) + 2)
}

func (mh *MinHeap) PrintHeap() {
	heapy := mh.heap

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
