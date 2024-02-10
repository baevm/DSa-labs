package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
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
	size    int

	stats Stats
}

const MAX_PROBE_COUNT = 40
const LIMIT_COEF = 0.5
const EMPTY_KEY = -1

var (
	ErrMaxProbe     = errors.New("ошибка: Достигнуто максимальное количество проб")
	ErrElemNotFound = errors.New("ошибка: элемент не найден")
)

func NewHTable(initialSize int) *HTable {
	stats := Stats{
		usedBucketsCount: 0,
		amountOfProbes:   0,
	}

	ht := &HTable{
		buckets: make([]HValue, initialSize),
		stats:   stats,
		size:    initialSize,
	}

	for i := range ht.buckets {
		ht.buckets[i].key = EMPTY_KEY
	}

	return ht
}

func (ht *HTable) Get(key int) (HValue, error) {
	hash := HashValue(key, len(ht.buckets))

	hash0 := hash
	probeCount := 0
	size := len(ht.buckets)

	if ht.buckets[hash].key == EMPTY_KEY {
		return HValue{}, ErrElemNotFound
	}

	for ht.buckets[hash].value != key && probeCount < MAX_PROBE_COUNT {
		hash = (hash0 + probeCount*probeCount) % size
		probeCount += 1
	}

	if probeCount >= MAX_PROBE_COUNT {
		return HValue{}, ErrMaxProbe
	}

	return ht.buckets[hash], nil
}

func (ht *HTable) Set(key, value int) (int, bool) {
	hash := HashValue(key, len(ht.buckets))

	hash0 := hash

	isSet := false
	probeCount := 0
	size := len(ht.buckets)

	if ht.buckets[hash].key == EMPTY_KEY {
		ht.buckets[hash] = HValue{key: hash, value: value}
		isSet = true
	} else {
		// Квадратичная проба
		for ht.buckets[hash].key != EMPTY_KEY && probeCount < MAX_PROBE_COUNT {
			hash = (hash0 + probeCount*probeCount) % size
			probeCount += 1

			ht.stats.amountOfProbes += 1
			isSet = true
		}

		// Линейная проба
		if !isSet && probeCount >= MAX_PROBE_COUNT {
			for ht.buckets[hash].key != EMPTY_KEY {
				hash = hash + 1
			}
		}

		ht.buckets[hash] = HValue{key: hash, value: value}
	}

	ht.stats.usedBucketsCount += 1

	if ht.getCoef() >= LIMIT_COEF {
		ht.ExtendTable()
	}

	return hash, isSet
}

func (ht *HTable) Delete(value int) (int, bool) {
	elem, err := ht.Get(value)

	if err != nil {
		return EMPTY_KEY, false
	}

	idx := elem.key

	ht.buckets[idx] = HValue{
		key:   EMPTY_KEY,
		value: 0,
	}

	return idx, true
}

func (ht *HTable) ExtendTable() {
	var newSize int

	if ht.size == 1 {
		newSize = 2
	} else {
		newSize = int(float64(ht.size) * 0.5)
	}

	for i := 0; i < newSize; i++ {
		ht.buckets = append(ht.buckets, HValue{key: EMPTY_KEY})
	}

	ht.size += newSize
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

func (ht *HTable) PrintStats() {
	coef := ht.getCoef()
	avgProbes := float64(ht.stats.amountOfProbes) / float64(ht.stats.usedBucketsCount)

	fmt.Printf("Число заполненных ячеек: %d\n", ht.stats.usedBucketsCount)
	fmt.Printf("Размер таблицы: %d\n", ht.size)
	fmt.Printf("Общее число проб: %d\n\n", ht.stats.amountOfProbes)
	fmt.Printf("Коэффициент заполнения: %f\n", coef)
	fmt.Printf("Среднее число проб: %f\n", avgProbes)
}

func (ht *HTable) AddRandomElements(min, max int) int {
	elements := getRandomElemets(min, max, ht.size/2)

	for _, v := range elements {
		_, isOk := ht.Set(v, v)

		if !isOk {
			fmt.Printf("Failed to set value:%d\n", v)
		}
	}

	return len(elements)
}

func (ht *HTable) getCoef() float64 {
	return float64(ht.stats.usedBucketsCount) / float64(ht.size)
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

	fmt.Printf("середина: %d\t hash значение: %d\t значение: %d\t квадрат: %d\n", middleDigits, hashValue, value, squared)

	return hashValue
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
