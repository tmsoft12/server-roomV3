package wsController

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/gofiber/websocket/v2"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	dataFile  = "sensor_data.json"
)

func Register(c *websocket.Conn) {
	clientsMu.Lock()
	clients[c] = true
	clientsMu.Unlock()

	err := sendJSONFileContents(c)
	if err != nil {
		log.Println("❌ Error Sending JSON File Contents:", err)
	}

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

type SensorData struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Topic   string `json:"topic"`
}

func sendJSONFileContents(c *websocket.Conn) error {
	fileData, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("⚠️ JSON File Does Not Exist:", dataFile)
			return nil
		}
		return err
	}

	if len(fileData) == 0 || string(fileData) == "{}" {
		log.Println("⚠️ JSON File is Empty")
		return nil
	}

	var jsonData map[string]SensorData
	if err := json.Unmarshal(fileData, &jsonData); err != nil {
		return err
	}

	for _, sensorData := range jsonData {
		flatJSON, err := json.Marshal(sensorData)
		if err != nil {
			log.Printf("❌ Error Marshalling Sensor Data for %s: %v", sensorData.Topic, err)
			continue
		}

		err = c.WriteMessage(websocket.TextMessage, flatJSON)
		if err != nil {
			log.Printf("❌ WebSocket Write Error for %s: %v", sensorData.Topic, err)
			return err
		}

	}

	return nil
}

func Broadcast(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("❌ WebSocket Write Error:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
