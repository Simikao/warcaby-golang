package game_handler

import (
	"net/http"
	"strconv"
	"sync"
	"warcaby/game"
	"warcaby/mqtt"

	"github.com/gin-gonic/gin"

	db "warcaby/database"
)

var (
	games   = make(map[int]*game.Game)
	gamesMu sync.Mutex
	nextID  = 1
)

func CreateGame(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak identyfikatora użytkownika"})
		return
	}
	userID := userIDVal.(int)

	gamesMu.Lock()
	defer gamesMu.Unlock()
	newGame := game.NewGame(nextID, userID)
	games[nextID] = newGame
	nextID++

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

	response := gin.H{
		"ID":            g.ID,
		"Board":         g.Board,
		"CurrentPlayer": g.CurrentPlayer,
		"Winner":        g.Winner,
		"Player1ID":     g.Player1ID,
		"Player2ID":     g.Player2ID,
	}

	if g.Player1ID != 0 {
		var user1 db.User
		if err := db.DB.First(&user1, g.Player1ID).Error; err == nil {
			response["Player1Nick"] = user1.Nick
		} else {
			response["Player1Nick"] = nil
		}
	} else {
		response["Player1Nick"] = nil
	}

	if g.Player2ID != 0 {
		var user2 db.User
		if err := db.DB.First(&user2, g.Player2ID).Error; err == nil {
			response["Player2Nick"] = user2.Nick
		} else {
			response["Player2Nick"] = nil
		}
	} else {
		response["Player2Nick"] = nil
	}
	c.JSON(http.StatusOK, response)
}

func MoveGame(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak identyfikatora użytkownika"})
		return
	}
	userID := userIDVal.(int)

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

	if err := g.Move(move.FromX, move.FromY, move.ToX, move.ToY, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mqtt.PublishGameUpdate(g)
	if g.Winner != 0 {
		mqtt.PublishGameWin(g.ID, int(g.Winner))
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
	defer gamesMu.Unlock()
	_, ok := games[id]
	if !ok {
		gamesMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Gra nie znaleziona"})
		return
	}
	delete(games, id)
	c.JSON(http.StatusOK, gin.H{"message": "Gra usunięta"})
}

func InviteUser(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak identyfikatora użytkownika"})
		return
	}
	userID := userIDVal.(int)

	idStr := c.Param("id")
	gameID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID gry"})
		return
	}

	gamesMu.Lock()
	g, ok := games[gameID]
	gamesMu.Unlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gra nie znaleziona"})
		return
	}

	if g.Player1ID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tylko twórca gry może zapraszać innych"})
		return
	}

	var payload struct {
		InviteeID int `json:"inviteeID"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Błąd dekodowania JSON"})
		return
	}
	if payload.InviteeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowy identyfikator zapraszanego użytkownika"})
		return
	}

	g.Player2ID = payload.InviteeID

	c.JSON(http.StatusOK, g)
}

func GetGames(c *gin.Context) {
	gamesMu.Lock()
	defer gamesMu.Unlock()

	var gameList []gin.H
	for _, g := range games {
		gameList = append(gameList, gin.H{
			"ID": g.ID,
		})
	}
	c.JSON(http.StatusOK, gameList)
}
