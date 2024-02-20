package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Сгенерировать 24 неповторяющихся трехзначных элементов.
// 1. Вывести их на экран.
// 2. Показать процесс построения Min-Heap-Tree.
// 3. Вывести на экран полученное дерево.
// Дополнительное задание: в отчете перечислить последовательность
// вершин построенного Min-Heap-Tree, соответствующую прямому порядку прохождения (NLR).

func main() {
	const size = 24
	const min = 99
	const max = 999

	randomElements := getRandomElemets(min, max, size)

	fmt.Println("Случайные элементы: ", randomElements)

	// minheap heap
	minheap := NewMinHeap()

	for _, elem := range randomElements {
		minheap.Push(Int(elem))
	}

	fmt.Println("MinHeap=", minheap)
	fmt.Println("size=", len(minheap))

	minheap.PrintHeap()

	fmt.Println("Traverse:")
	minheap.Traverse(0)
}

func getRandomElemets(min, max, size int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// массив из случайных неповторяющихся чисел
	p := r.Perm(max - min + 1)

	// добавляем к каждому минимальное число
	// чтобы избежать чисел меньше минимально заданного
	for i := range p {
		p[i] += min
	}

	return p[:size]
}
