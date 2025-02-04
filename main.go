package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	db "warcaby/database"
	game "warcaby/game"
	gHandler "warcaby/handlers/game_handler"
	uHandler "warcaby/handlers/user_handler"
	mqtt "warcaby/mqtt"

	"github.com/gin-gonic/gin"
)

func main() {
	debug := os.Getenv("LOGFILE")
	if debug != "" {
		f, _ := os.Create(debug + ".log")
		defer f.Close()
		gin.DefaultWriter = io.MultiWriter(f)
	}

	db.InitDB()
	mqtt.InitMQTT()
	r := gin.Default()

	authorized := r.Group("/", AuthMiddleware())
	{
		authorized.POST("/games/new", gHandler.CreateGame)
		authorized.POST("/games/:id/invite", gHandler.InviteUser)
		authorized.PUT("/games/:id/move", gHandler.MoveGame)
		authorized.DELETE("/games/:id", gHandler.DeleteGame)

		authorized.GET("/users/me", uHandler.GetMyUser)
		authorized.PUT("/users/:id", uHandler.UpdateUser)
		authorized.DELETE("/users/:id", uHandler.DeleteUser)
	}

	r.GET("/games/:id", gHandler.GetGame)
	r.GET("/games/list", gHandler.GetGames)

	r.GET("/users/:id", uHandler.GetUser)
	r.GET("/users", uHandler.SearchUsers)

	r.POST("/register", uHandler.RegisterUser)
	r.POST("/login", uHandler.Login)

	r.StaticFile("/game", "./game.html")
	r.StaticFile("/login", "./login.html")
	r.StaticFile("/search", "./search.html")
	r.StaticFile("/profile", "./profile.html")

	if err := r.RunTLS(":8080", "./cert/cert.pem", "./cert/key.pem"); err != nil {
		log.Fatal(err)
	}

	g := game.NewGame(1, 1)
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

		err := g.Move(fromX, fromY, toX, toY, 1)
		if err != nil {
			fmt.Println("Błąd:", err)
		}
	}

	fmt.Println("Koniec gry!")
}
