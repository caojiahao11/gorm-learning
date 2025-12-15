package main

import "fmt"

func main() {
	func() {
		fmt.Print("这是一个匿名函数")
	}()

	sayHello := func(name string) {
		fmt.Printf("Hello, %s!\n", name)
	}
	sayHello("Alice")
	sayHello("Bob")
	sayHello("Charlie")

	executeOperation := func(a, b int, operation func(int, int) int) int {
		return operation(a, b)
	}
	sum1 := executeOperation(5, 3, func(a, b int) int {
		return a + b
	})
	fmt.Println(sum1)

	executeOperation(6, 1, func(a, b int) int {
		return a - b
	})
	executeOperation(1, 2, func(a, b int) int {
		return a * b
	})

}
