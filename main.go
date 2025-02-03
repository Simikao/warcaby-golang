package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	db "warcaby/database"
	game "warcaby/game"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	games   = make(map[int]*game.Game)
	gamesMu sync.Mutex
	nextID  = 1
)

func createGame(c *gin.Context) {
	gamesMu.Lock()
	newGame := game.NewGame(nextID)
	games[nextID] = newGame
	nextID++
	gamesMu.Unlock()

	c.JSON(http.StatusOK, newGame)
}

func getGame(c *gin.Context) {
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

func moveGame(c *gin.Context) {
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

func deleteGame(c *gin.Context) {
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

func registerUser(c *gin.Context) {
	var newUser db.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Błąd dekodowania JSON"})
		return
	}

	if newUser.Nick == "" || newUser.Email == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nick, email i hasło są wymagane"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy haszowaniu hasła"})
		return
	}
	newUser.Password = string(hashedPassword)

	if err := db.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas tworzenia użytkownika"})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

func login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Błąd dekodowania JSON"})
		return
	}

	var user db.User
	if err := db.DB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Niepoprawny email lub hasło"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Niepoprawny email lub hasło"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logowanie udane", "user": user})
}

func main() {
	db.InitDB()
	r := gin.Default()

	r.POST("/games/new", createGame)
	r.GET("/games/:id", getGame)
	r.PUT("/games/:id/move", moveGame)
	r.DELETE("/games/:id", deleteGame)

	r.POST("/register", registerUser)
	r.POST("/login", login)

	r.StaticFile("/game", "./game.html")

	r.Run(":8080")

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
