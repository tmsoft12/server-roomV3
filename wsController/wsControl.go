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
	dataFile  = "sensor_data.json" // Path to your JSON file
)

func Register(c *websocket.Conn) {
	clientsMu.Lock()
	clients[c] = true
	clientsMu.Unlock()

	// Send all contents of the JSON file to the newly connected client in flat format
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

	// Keep the connection alive to listen for incoming messages
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

// SensorData represents the flat structure of the data to be sent over WebSocket
type SensorData struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Topic   string `json:"topic"`
}

// sendJSONFileContents reads the JSON file and sends all entries in flat format
func sendJSONFileContents(c *websocket.Conn) error {
	// Read the JSON file
	fileData, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("⚠️ JSON File Does Not Exist:", dataFile)
			return nil // No file exists yet, no error to propagate
		}
		return err
	}

	// Check if the file is empty
	if len(fileData) == 0 || string(fileData) == "{}" {
		log.Println("⚠️ JSON File is Empty")
		return nil // Nothing to send
	}

	// Unmarshal the JSON file into a map
	var jsonData map[string]SensorData
	if err := json.Unmarshal(fileData, &jsonData); err != nil {
		return err
	}

	// Send each entry as a separate flat JSON message
	for _, sensorData := range jsonData {
		flatJSON, err := json.Marshal(sensorData)
		if err != nil {
			log.Printf("❌ Error Marshalling Sensor Data for %s: %v", sensorData.Topic, err)
			continue
		}

		// Send the flat JSON data to the client
		err = c.WriteMessage(websocket.TextMessage, flatJSON)
		if err != nil {
			log.Printf("❌ WebSocket Write Error for %s: %v", sensorData.Topic, err)
			return err // Return error if the client disconnects
		}

		log.Printf("✅ Sent Flat JSON Data to WebSocket Client: %s", string(flatJSON))
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
