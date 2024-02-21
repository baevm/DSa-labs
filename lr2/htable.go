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
const EMPTY_VAL = -1

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
		ht.buckets[i].value = EMPTY_VAL
		ht.buckets[i].key = i
	}

	return ht
}

func (ht *HTable) Get(key int) (HValue, error) {
	hash := HashValue(key, len(ht.buckets))

	hash0 := hash
	probeCount := 0
	size := len(ht.buckets)

	if ht.buckets[hash].value == EMPTY_VAL {
		return HValue{}, ErrElemNotFound
	}

	// квадратичное пробирование
	for ht.buckets[hash].value != key && probeCount < MAX_PROBE_COUNT && hash < size {
		hash = (hash0 + probeCount*probeCount) % size
		probeCount += 1
	}

	if ht.buckets[hash].value == key {
		return ht.buckets[hash], nil
	}

	// линейное пробирование
	if probeCount >= MAX_PROBE_COUNT && ht.buckets[hash].value != key {
		hash = 0

		for hash < ht.size && ht.buckets[hash].value != EMPTY_VAL {

			hash += 1
			probeCount += 1
		}
	}

	if ht.buckets[hash].value == key {
		return ht.buckets[hash], nil
	}

	return HValue{}, ErrElemNotFound
}

func (ht *HTable) Set(key, value int) (int, bool) {
	hash := HashValue(key, len(ht.buckets))

	hash0 := hash

	isSet := false
	probeCount := 0
	size := len(ht.buckets)

	if ht.buckets[hash].value == EMPTY_VAL {
		ht.buckets[hash] = HValue{key: hash, value: value}
		isSet = true

		ht.stats.amountOfProbes += 1
	} else {
		// Квадратичная проба
		for ht.buckets[hash].value != EMPTY_VAL && probeCount < MAX_PROBE_COUNT {
			hash = (hash0 + probeCount*probeCount) % size
			probeCount += 1

			// fmt.Println("quad hash: ", hash, "probeCount: ", probeCount)

			ht.stats.amountOfProbes += 1
		}

		if ht.buckets[hash].value == EMPTY_VAL && probeCount <= MAX_PROBE_COUNT {
			isSet = true
			ht.buckets[hash] = HValue{key: hash, value: value}
		}
	}

	// Линейная проба
	if !isSet && probeCount >= MAX_PROBE_COUNT {
		hash = 0

		for hash < ht.size && ht.buckets[hash].value != EMPTY_VAL {
			hash = hash + 1

			// fmt.Println("linear hash: ", hash, "probeCount: ", probeCount)

			ht.stats.amountOfProbes += 1
		}

		if ht.buckets[hash].value == EMPTY_VAL {
			isSet = true
			ht.buckets[hash] = HValue{key: hash, value: value}
		}
	}

	ht.stats.usedBucketsCount += 1

	// if ht.getCoef() >= LIMIT_COEF {
	// 	ht.ExtendTable()
	// }

	return hash, isSet
}

func (ht *HTable) Delete(value int) (int, bool) {
	elem, err := ht.Get(value)

	if err != nil {
		return EMPTY_VAL, false
	}

	idx := elem.key

	// TODO: пустое удаленное значение должно быть с другим value. например -2
	// при поиске проверять если -2, то значит оно тоже пустое
	ht.buckets[idx] = HValue{
		key:   idx,
		value: EMPTY_VAL,
	}

	ht.stats.usedBucketsCount -= 1

	return idx, true
}

func (ht *HTable) Change(oldElement, newElement int) (int, int, error) {
	newElemExist, _ := ht.Get(newElement)

	if newElemExist.value == newElement {
		return 0, 0, fmt.Errorf("Элемент %d уже существует. \n\n", newElement)
	}

	oldIdx, isDeleted := ht.Delete(oldElement)

	if !isDeleted {
		return 0, 0, fmt.Errorf("Элемент %d не найден. \n\n", oldElement)
	}

	idx, isSet := ht.Set(newElement, newElement)

	if !isSet {
		return 0, 0, fmt.Errorf("Что то пошло не так... Элемент: %d \n\n", newElement)
	}

	return idx, oldIdx, nil
}

func (ht *HTable) ExtendTable() {
	var newSize int

	if ht.size == 1 {
		newSize = 2
	} else {
		newSize = int(float64(ht.size) * 0.5)
	}

	for i := 0; i < newSize; i++ {
		// fix
		ht.buckets = append(ht.buckets, HValue{key: ht.size + i, value: EMPTY_VAL})
	}

	ht.size += newSize
}

func (ht *HTable) Print() {
	// fmt.Println("|-------|-------|")
	// fmt.Printf("| %-5s | %-5s |\n", "key", "value")
	// fmt.Println("|-------|-------|")
	// for _, h := range ht.buckets {
	// 	fmt.Printf("| %-5d | %-5d |\n", h.key, h.value)
	// }
	fmt.Printf("%+v \n\n", ht.buckets[0:10])
	fmt.Printf("%+v \n\n", ht.buckets[10:20])
	fmt.Printf("%+v \n\n", ht.buckets[20:30])
	fmt.Printf("%+v \n\n", ht.buckets[30:40])
	fmt.Printf("%+v \n\n", ht.buckets[40:50])
	fmt.Printf("%+v \n\n", ht.buckets[50:60])
	fmt.Printf("%+v \n\n", ht.buckets[70:])

	// fmt.Println("|-------|-------|")
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

func (ht *HTable) AddRandomElements(min, max, size int, isPrintValue bool) int {
	elements := getRandomElemets(min, max, size)

	for _, v := range elements {
		// elem, _ := ht.Get(v)

		// if elem.value != -1 {
		// 	return 0
		// }

		idx, isOk := ht.Set(v, v)

		if !isOk {
			fmt.Printf("Failed to set value:%d\n", v)
		}

		if isPrintValue {
			fmt.Printf("Индекс добавленного случайного элемента %v: %v \n", v, idx)
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

	// кол-во цифр из середины
	digits := 2

	start := int(math.Log10(float64(squared))) / 2

	// достаем середину
	middleDigits := (squared / int(math.Pow10(start))) % int(math.Pow10(digits))

	// модуль от размера таблицы чтобы получить валидный ключ
	hashValue := middleDigits % tableSize

	// fmt.Printf("середина: %d\t hash значение: %d\t значение: %d\t квадрат: %d\t\n ", middleDigits, hashValue, value, squared)

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

	// fmt.Println("min ", min, "max ", max, "size ", size, "p ", p, "p len", len(p))

	return p[:size]
}
