package game

import (
	"errors"
	"fmt"
)

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
	fmt.Println("Ruch gracza:", pieceToString(g.CurrentPlayer))
}

func pieceToString(p Piece) string {
	switch p {
	case Black:
		return "Czarne"
	case White:
		return "Białe"
	default:
		return "Brak"
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *Game) Move(fromX, fromY, toX, toY int) error {
	if fromX < 0 || fromX >= BoardSize || fromY < 0 || fromY >= BoardSize ||
		toX < 0 || toX >= BoardSize || toY < 0 || toY >= BoardSize {
		return errors.New("współrzędne poza zakresem planszy")
	}

	piece := g.Board[fromX][fromY]
	if piece == Empty {
		return errors.New("na wskazanej pozycji nie ma pionka")
	}

	if piece != g.CurrentPlayer {
		return errors.New("nie możesz poruszać pionkiem przeciwnika")
	}

	dx := toX - fromX
	dy := toY - fromY

	if abs(dx) != abs(dy) {
		return errors.New("ruch musi być po przekątnej")
	}

	if abs(dx) == 1 {
		if piece == White && dx != -1 {
			return errors.New("białe pionki mogą poruszać się tylko do góry")
		}
		if piece == Black && dx != 1 {
			return errors.New("czarne pionki mogą poruszać się tylko w dół")
		}
		if g.Board[toX][toY] != Empty {
			return errors.New("docelowe pole nie jest puste")
		}

		g.Board[toX][toY] = piece
		g.Board[fromX][fromY] = Empty

	} else if abs(dx) == 2 {
		if piece == White && dx != -2 {
			return errors.New("białe pionki mogą bić tylko do góry")
		}
		if piece == Black && dx != 2 {
			return errors.New("czarne pionki mogą bić tylko w dół")
		}

		midX := fromX + dx/2
		midY := fromY + dy/2

		opponent := Black
		if piece == Black {
			opponent = White
		}

		if g.Board[midX][midY] != opponent {
			return errors.New("brak pionka przeciwnika do zbicia")
		}

		if g.Board[toX][toY] != Empty {
			return errors.New("docelowe pole nie jest puste")
		}

		g.Board[toX][toY] = piece
		g.Board[fromX][fromY] = Empty
		g.Board[midX][midY] = Empty
	} else {
		return errors.New("nieprawidłowy ruch: ruch musi być o jedno pole")
	}

	if g.CurrentPlayer == White {
		g.CurrentPlayer = Black
	} else {
		g.CurrentPlayer = White
	}

	return nil
}
