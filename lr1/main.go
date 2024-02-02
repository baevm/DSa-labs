package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Вариант - 2
// хэш таблица из элементов m = 49
// размерность элемента n = 2
// хэш функция - средняя часть квадрата ключа
// метод разрешения конфилкта - квадратичные пробы

func main() {
	const size = 49
	const min = 10
	const max = 99

	randomNumbers := getRandomElemets(min, max, size)

	fmt.Println("Случайные числа: ", randomNumbers)

	hTableSize := int(math.Floor(size * 1.5))
	hashTable := NewHTable(hTableSize)

	for _, v := range randomNumbers {
		isOk := hashTable.Set(v, v)

		if !isOk {
			fmt.Printf("Failed to set value:%d\n", v)
		}
	}

	hashTable.Print()

	coef := float64(hashTable.stats.usedBucketsCount) / float64(hTableSize)
	avgProbes := float64(hashTable.stats.amountOfProbes) / float64(hashTable.stats.usedBucketsCount)

	fmt.Printf("Число заполненных ячеек: %d\n", hashTable.stats.usedBucketsCount)
	fmt.Printf("Общее число проб: %d\n\n", hashTable.stats.amountOfProbes)
	fmt.Printf("Коэффициент заполнения: %f\n", coef)
	fmt.Printf("Среднее число проб: %f\n", avgProbes)
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
