package wsController

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

func Register(c *websocket.Conn) {
	clientsMu.Lock()
	clients[c] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, c)
		clientsMu.Unlock()
		c.Close()
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

func Broadcast(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("WebSocket Write Error:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
