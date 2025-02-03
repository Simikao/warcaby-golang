package game_handler

import (
	"net/http"
	"strconv"
	"sync"
	"warcaby/game"

	"github.com/gin-gonic/gin"
)

var (
	games   = make(map[int]*game.Game)
	gamesMu sync.Mutex
	nextID  = 1
)

func CreateGame(c *gin.Context) {
	gamesMu.Lock()
	newGame := game.NewGame(nextID)
	games[nextID] = newGame
	nextID++
	gamesMu.Unlock()

	c.JSON(http.StatusOK, newGame)
}

func GetGame(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID"})
		return
	}

	gamesMu.Lock()
	g, ok := games[id]
	gamesMu.Unlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gra nie znaleziona"})
		return
	}

	c.JSON(http.StatusOK, g)
}

func MoveGame(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID"})
		return
	}

	gamesMu.Lock()
	g, ok := games[id]
	gamesMu.Unlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gra nie znaleziona"})
		return
	}

	var move struct {
		FromX int `json:"fromX"`
		FromY int `json:"fromY"`
		ToX   int `json:"toX"`
		ToY   int `json:"toY"`
	}
	if err := c.ShouldBindJSON(&move); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Błąd dekodowania JSON"})
		return
	}

	if err := g.Move(move.FromX, move.FromY, move.ToX, move.ToY); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, g)
}

func DeleteGame(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID"})
		return
	}

	gamesMu.Lock()
	_, ok := games[id]
	if !ok {
		gamesMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Gra nie znaleziona"})
		return
	}
	delete(games, id)
	gamesMu.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "Gra usunięta"})
}
