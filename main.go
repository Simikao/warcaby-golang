package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	game "warcaby/warcaby"
)

func main() {
	g := game.NewGame(1)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		g.PrintBoard()

		fmt.Println("Podaj ruch w formacie: fromX fromY toX toY (lub wpisz 'exit' aby zakończyć):")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "exit" {
			break
		}

		parts := strings.Fields(input)
		if len(parts) != 4 {
			fmt.Println("Nieprawidłowy format. Spróbuj ponownie.")
			continue
		}

		fmt.Println(parts)
	}
}
