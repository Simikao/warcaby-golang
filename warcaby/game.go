package game

import "fmt"

const BoardSize = 8

type Piece int

const (
	Empty Piece = iota
	Black
	White
)

type Game struct {
	ID            int
	Board         [BoardSize][BoardSize]Piece
	CurrentPlayer Piece
	Winner        Piece
}

func NewGame(id int) *Game {
	game := &Game{
		ID:            id,
		CurrentPlayer: White,
		Winner:        Empty,
	}

	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			if (i+j)%2 != 0 {
				if i < 3 {
					game.Board[i][j] = Black
				} else if i >= BoardSize-3 {
					game.Board[i][j] = White
				} else {
					game.Board[i][j] = Empty
				}
			} else {
				game.Board[i][j] = Empty
			}
		}
	}

	return game
}

func (g *Game) PrintBoard() {
	fmt.Print("  ")
	for j := 0; j < BoardSize; j++ {
		fmt.Printf("%d ", j)
	}
	fmt.Println()

	for i := 0; i < BoardSize; i++ {
		fmt.Printf("%d ", i)

		for j := 0; j < BoardSize; j++ {
			var ch string
			switch g.Board[i][j] {
			case Black:
				ch = "B"
			case White:
				ch = "W"
			default:
				ch = "."
			}
			fmt.Print(ch, " ")
		}
		fmt.Println()
	}
}
