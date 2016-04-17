package main

import "fmt"

const (
	CLEAR = "\033[H\033[2J"
)

func ConsoleClear() {
	fmt.Print(CLEAR)
}
