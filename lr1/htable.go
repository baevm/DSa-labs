package main

import (
	"fmt"
	"math"
)

type Stats struct {
	usedBucketsCount int // Количество заполненных ячеек в таблице
	amountOfProbes   int // Количество проб
}

type HValue struct {
	key   int
	value int
}

type HTable struct {
	buckets []HValue // Ячейки хэш таблицы

	stats Stats
}

func NewHTable(initialSize int) *HTable {
	stats := Stats{
		usedBucketsCount: 0,
		amountOfProbes:   0,
	}

	ht := &HTable{
		buckets: make([]HValue, initialSize),
		stats:   stats,
	}

	for i := range ht.buckets {
		ht.buckets[i].key = -1
	}

	return ht
}

func (ht *HTable) Get(key int) (HValue, bool) {
	hash := HashValue(key, len(ht.buckets))

	hash0 := hash
	probeCount := 0
	size := len(ht.buckets)

	if ht.buckets[hash].key == -1 {
		return HValue{}, false
	}

	for ht.buckets[hash].value != key && probeCount < 30 {
		hash = (hash0 + probeCount*probeCount) % size
		probeCount += 1
	}

	return ht.buckets[hash], false
}

func (ht *HTable) Set(key, value int) bool {
	hash := HashValue(key, len(ht.buckets))

	hash0 := hash

	isSet := false
	probeCount := 0
	size := len(ht.buckets)

	if ht.buckets[hash].key == -1 {
		ht.buckets[hash] = HValue{key: hash, value: value}
		isSet = true
	} else {
		// Квадратичная проба
		for ht.buckets[hash].key != -1 && probeCount < 30 {
			hash = (hash0 + probeCount*probeCount) % size
			probeCount += 1

			ht.stats.amountOfProbes += 1
			isSet = true
		}

		// Линейная проба
		if !isSet && probeCount >= 30 {
			for ht.buckets[hash].key != -1 {
				hash = hash + 1
			}
		}

		ht.buckets[hash] = HValue{key: hash, value: value}
	}

	ht.stats.usedBucketsCount += 1

	return isSet
}

func (ht *HTable) Print() {
	fmt.Println("|-------|-------|")
	fmt.Printf("| %-5s | %-5s |\n", "key", "value")
	fmt.Println("|-------|-------|")
	for _, h := range ht.buckets {
		fmt.Printf("| %-5d | %-5d |\n", h.key, h.value)
	}
	fmt.Println("|-------|-------|")
}

// хэш функция - средняя часть квадрата ключа
func HashValue(value, tableSize int) int {
	squared := value * value

	// длина квадрата числа
	length := int(math.Log10(float64(squared))) + 1

	// кол-во цифр из середины
	digits := 0

	if length%2 == 0 {
		digits = 2
	} else {
		digits = 1
	}

	start := int(math.Log10(float64(squared))) / 2

	// достаем середину
	middleDigits := (squared / int(math.Pow10(start))) % int(math.Pow10(digits))

	// модуль от размера таблицы чтобы получить валидный ключ
	hashValue := middleDigits % tableSize

	fmt.Printf("середина: %d\t hash значение: %d\t значение: %d\t квадрат: %d\t\n ", middleDigits, hashValue, value, squared)

	return hashValue
}
