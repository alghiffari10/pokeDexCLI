package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := scanner.Text()
		lowerText := strings.ToLower(input)
		words := cleanInput(lowerText)
		fmt.Printf("Your command was: %s\n", words[0])

	}

}
