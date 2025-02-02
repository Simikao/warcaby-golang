package game

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
