package main

import "fmt"

func main() {
	r := Add(10, 20)
	fmt.Println(r)
}

func Add(a, b int) int {
	return a + b
}
