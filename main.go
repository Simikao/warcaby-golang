package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

		fromX, err1 := strconv.Atoi(parts[0])
		fromY, err2 := strconv.Atoi(parts[1])
		toX, err3 := strconv.Atoi(parts[2])
		toY, err4 := strconv.Atoi(parts[3])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Println("Wprowadź poprawne liczby!")
			continue
		}

		err := g.Move(fromX, fromY, toX, toY)
		if err != nil {
			fmt.Println("Błąd:", err)
		}
	}

	fmt.Println("Koniec gry!")
}
