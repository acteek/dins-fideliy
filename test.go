package main

import "fmt"

type Test struct {
	Field string
}

const (
	lol = "looool"
)

func main() {
	a := Test{lol}

	a.Field = "string1"
	fmt.Print(a.Field)

}
