package main

import (
	"fmt"
	"os"
)

func main() {
	name := os.Getenv("NAME")

	fmt.Printf("Hello %s from the Console", name)
}
