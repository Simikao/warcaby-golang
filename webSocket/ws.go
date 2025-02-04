package webSocket

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"warcaby/game"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	clients   = make(map[int]map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

func WsGameHandler(c *gin.Context) {
	fmt.Println("I'm here")
	idStr := c.Param("id")
	gameID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID gry"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Błąd przy upgrade'owaniu połączenia:", err)
		return
	}

	clientsMu.Lock()
	if clients[gameID] == nil {
		clients[gameID] = make(map[*websocket.Conn]bool)
	}
	clients[gameID][conn] = true
	clientsMu.Unlock()

	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
		clientsMu.Lock()
		delete(clients[gameID], conn)
		clientsMu.Unlock()
	}()
	i := 0
	for {
		select {
		case <-ticker.C:
			i++
			message := "Aktualizacja gry #" + strconv.Itoa(i)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Błąd przy wysyłaniu wiadomości przez WebSocket:", err)
				return
			}
		}
	}
}

func BroadcastGameUpdate(g *game.Game) {
	msg := map[string]interface{}{
		"type": "move",
		"data": g,
	}
	clientsMu.Lock()
	conns, exists := clients[g.ID]
	clientsMu.Unlock()
	if !exists {
		return
	}
	for conn := range conns {
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Błąd wysyłania JSON przez WebSocket:", err)
			conn.Close()
			clientsMu.Lock()
			delete(conns, conn)
			clientsMu.Unlock()
		}
	}
}
