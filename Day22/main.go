package main

import "fmt"

func main() {
	f()
	fmt.Println("Завершение работы")
}

func f() {
	defer func() {
		if r := recover(); r != nil { // если случилась паника
			fmt.Println("Восстановление со значением", r)
		}
	}()
	fmt.Println("вызываем g")
	g(1)
	fmt.Println("завершение работы g")
}

func g(i int) {
	if i > 3 { // сделаем искусственную панику
		fmt.Println("Паника!")
		panic(fmt.Sprintf("%v", i))
	}
	defer fmt.Println(i, "defer из g") // заполним стек вызовов
	fmt.Println(i, "вызов из g")
	g(i + 1) // снова запустим g из самой себя
}
