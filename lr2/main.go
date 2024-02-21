package main

import (
	"errors"
	"fmt"
	"math"
	"os"
)

// Добавить функции позволяющие:
// 1. Сгенерировать или ввести в интерактивном режиме новые элементы.
// 2. Осуществить поиск элемента в таблице.
// 3. Добавить элемент в таблицу.
// 4. Удалить элемент из таблицы.
// 5. Ввести в интерактивном режиме:
// 6. Вывести статистику:
//    коэффициент заполнения таблицы
//	  среднее число шагов необходимых для размещения некоторого ключа в таблице.

const (
	ADD = iota + 1
	FIND
	DEL
	CHANGE
	STATS
	PRINT
)

func main() {
	const size = 49
	const min = 10
	const max = 99

	hTableSize := int(math.Floor(size * 1.5))

	hashTable := NewHTable(hTableSize)

	hashTable.AddRandomElements(min, max, size, false)

	fmt.Println("Начальная таблица:")
	hashTable.Print()

	for {
		fmt.Printf(`
Выберите нужный вариант:
%d. Добавить элемент
%d. Поиск элемента
%d. Удалить элемент
%d. Заменить элемент
%d. Вывести статистику (коэффициент заполнения таблицы, среднее число шагов необходимых для размещения ключа в таблице)
%d. Вывести таблицу
`, ADD, FIND, DEL, CHANGE, STATS, PRINT)

		fmt.Print("Ваш выбор: ")

		var answer int
		_, err := fmt.Scanf("%d", &answer)

		if err != nil {
			fmt.Println("Что то пошло не так...")
			os.Exit(1)
		}

		switch answer {
		case ADD:
			fmt.Println("Выберите вариант:")
			fmt.Println("1. Добавить в интерактивном режиме")
			fmt.Println("2. Сгенерировать случайные числа")

			var answerAdd int
			fmt.Scan(&answerAdd)

			switch answerAdd {
			case 1:
				var element int

				fmt.Print("Введите элемент: ")
				fmt.Scanf("%d", &element)

				if element < 1 {
					fmt.Printf("Элемент должен быть больше 0 \n\n")
					break
				}

				if hashTable.size == hashTable.stats.usedBucketsCount {
					fmt.Printf("Таблица полностью заполнена. Удалите элементы перед добавлением нового \n\n")
					break
				}

				elemExist, _ := hashTable.Get(element)

				if elemExist.value == element {
					fmt.Printf("Элемент уже существует. Индекс: %d \n\n", elemExist.key)
					break
				}

				idx, isSet := hashTable.Set(element, element)

				if !isSet {
					fmt.Printf("Что то пошло не так... Элемент: %d \n\n", element)
					break
				}

				fmt.Printf("Индекс нового элемента: %d \n\n", idx)

			case 2:
				var answerSize int
				fmt.Println("Введите размер:")

				fmt.Scan(&answerSize)

				for answerSize > 0 {
					addedLen := hashTable.AddRandomElements(min, max, 1, true)

					answerSize -= addedLen
				}

			default:
				break
			}

		case FIND:
			var element int

			fmt.Print("Введите элемент: ")
			fmt.Scanf("%d", &element)

			elem, err := hashTable.Get(element)

			if err != nil {
				if errors.Is(err, ErrElemNotFound) {
					fmt.Printf("Элемент %d не найден. \n\n", element)
					break
				}

				if errors.Is(err, ErrMaxProbe) {
					fmt.Printf("Достигнут максимум доступных проб \n\n")
					break
				}
			}

			fmt.Printf("Элемент %d находится на индексе: %d \n\n", element, elem.key)

		case DEL:
			var element int

			fmt.Print("Введите элемент: ")
			fmt.Scanf("%d", &element)

			idx, isDeleted := hashTable.Delete(element)

			if !isDeleted {
				fmt.Printf("Элемент %d не найден. \n\n", element)
				break
			}

			fmt.Printf("Элемент удален с индекса: %d \n\n", idx)

		case CHANGE:
			var element int

			fmt.Print("Введите заменяемый элемент: ")
			fmt.Scanf("%d", &element)

			if element < 1 {
				fmt.Printf("Элемент должен быть больше 0 \n\n")
				break
			}

			if hashTable.size == hashTable.stats.usedBucketsCount {
				fmt.Printf("Таблица полностью заполнена. Удалите элементы перед добавлением нового \n\n")
				break
			}

			_, err := hashTable.Get(element)

			if err != nil {
				fmt.Printf("Элемент %d не найден. \n\n", element)
				break
			}

			var newElement int

			fmt.Print("Введите новый элемент: ")
			fmt.Scanf("%d", &newElement)

			idx, oldIdx, err := hashTable.Change(element, newElement)

			if err != nil {
				fmt.Printf("%s \n\n", err.Error())
				break
			}

			fmt.Printf("Индекс удаленного элемента: %d \n\n", oldIdx)
			fmt.Printf("Индекс нового элемента: %d \n\n", idx)

		case STATS:
			fmt.Printf("Статистика:\n\n")
			hashTable.PrintStats()

		case PRINT:
			fmt.Printf("Таблица:\n\n")
			hashTable.Print()

		default:
			fmt.Printf("Такой опции не существует.")
		}
	}
}
