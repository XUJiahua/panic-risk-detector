package main

import "fmt"
import "strings"

func main() {
	hello := "Hello"
	world := "World"
	words := []string{hello, world}
	SayHello(words)
}

// SayHello says Hello
func SayHello(words []string) {
	fmt.Println(joinStrings(words))
}

// joinStrings joins strings
func joinStrings(words []string) string {
	go sample()
	return strings.Join(words, ", ")
}

func sample() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%+v\n", r)
		}
	}()

	// something happen
	panic("hello world")
	// something should happen
}
