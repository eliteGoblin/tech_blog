package main

import (
	"fmt"
	"utils"
)

func main() {
	fmt.Println("Hello", utils.Div(1, 2))
	utils.Div(1, 0)
	n := 0.0
	n = n + 1
	n = n - 1
	m := 1 / n
	n = m * 0
	fmt.Println(n)
}