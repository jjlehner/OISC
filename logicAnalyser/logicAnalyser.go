package main

import "fmt"

func main() {
	inputLines := "((!a||b&&c)&&d)"
	root, _ := descendExpression([]rune(inputLines), 0)
	a := root.eval(&map[string]bool{
		"a": false,
		"b": false,
		"c": true,
		"d": true,
	})
	fmt.Println(a)
}
