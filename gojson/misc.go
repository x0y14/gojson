package gojson

import "fmt"

func ShowPos(text string) {
	data := []rune(text)
	for i, r := range data {
		fmt.Printf("[%v] %v\n", i, string(r))
	}
}
