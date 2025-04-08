package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"server/controllers"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var client mqtt.Client
var topics map[string]string
var dataFile = "sensor_data.json"
var mu sync.Mutex

func ConnectMQTT() {
	topics = map[string]string{
		"temperature": os.Getenv("MQTT_TOPIC_TEMPERATURE"),
		"humidity":    os.Getenv("MQTT_TOPIC_HUMIDITY"),
		"main":        os.Getenv("MQTT_TOPIC_MAIN"),
		"door1":       os.Getenv("MQTT_TOPIC_1"),
		"door2":       os.Getenv("MQTT_TOPIC_2"),
		"door3":       os.Getenv("MQTT_TOPIC_3"),
		"door4":       os.Getenv("MQTT_TOPIC_4"),
		"door5":       os.Getenv("MQTT_TOPIC_5"),
		"door6":       os.Getenv("MQTT_TOPIC_6"),
		"door7":       os.Getenv("MQTT_TOPIC_7"),
		"door8":       os.Getenv("MQTT_TOPIC_8"),
		"door9":       os.Getenv("MQTT_TOPIC_9"),
	}

	broker := os.Getenv("MQTT_BROKER")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetDefaultPublishHandler(messageHandler)

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ MQTT Connection Error: %v", token.Error())
	}

	log.Println("✅ Connected to MQTT Broker:", broker)

	for friendlyName, topic := range topics {
		if topic == "" {
			log.Println("⚠️ Skipping Empty Topic:", friendlyName)
			continue
		}
		subscribe(topic)
	}
}

func subscribe(topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("❌ Subscription Error: %v", token.Error())
	}
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("📩 Received Message: %s, Topic: %s\n", msg.Payload(), msg.Topic())
	friendlyName := ""
	for name, topic := range topics {
		if topic == msg.Topic() {
			friendlyName = name
			break
		}
	}
	if friendlyName == "" {
		log.Printf("⚠️ No Friendly Name Found for Topic: %s", msg.Topic())
		return
	}

	// Prepare the sensor data in the specified format
	sensorData := map[string]string{
		"data":    string(msg.Payload()),
		"message": fmt.Sprintf("🚪 %s Sensor Triggered: %s", friendlyName, string(msg.Payload())),
		"topic":   friendlyName,
	}

	// Marshal the sensor data to JSON format
	dataJson, err := json.Marshal(sensorData)
	if err != nil {
		log.Printf("❌ Error Marshalling Sensor Data: %v", err)
		return
	}

	// Save or update the sensor data in the JSON file
	saveOrUpdateData(friendlyName, dataJson)
	controllers.HandleSensorData(friendlyName, string(msg.Payload()))
}

// Function to save or update the JSON file
func saveOrUpdateData(topic string, data []byte) {
	mu.Lock()
	defer mu.Unlock()

	// Check if the file exists and if not, create it with empty data
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		// If the file doesn't exist, create an empty JSON object
		initialData := "{}"
		err := os.WriteFile(dataFile, []byte(initialData), 0644) // Changed to os.WriteFile
		if err != nil {
			log.Printf("❌ Error Creating Initial JSON File: %v", err)
			return
		}
	}

	// Read the existing data from the file
	fileData, err := os.ReadFile(dataFile) // Changed to os.ReadFile
	if err != nil {
		log.Printf("❌ Error Reading File: %v", err)
		return
	}

	var existingData map[string]json.RawMessage
	if len(fileData) == 0 {
		// File is empty, initialize with an empty map
		existingData = make(map[string]json.RawMessage)
	} else {
		// Try unmarshalling the file data if it's not empty
		err = json.Unmarshal(fileData, &existingData)
		if err != nil {
			log.Printf("❌ Error Unmarshalling File Data: %v", err)
			// Fallback to an empty map if unmarshalling fails
			existingData = make(map[string]json.RawMessage)
		}
	}

	// Update or add the new sensor data
	existingData[topic] = data

	// Marshal the updated content back to JSON
	updatedData, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		log.Printf("❌ Error Marshalling Updated Data: %v", err)
		return
	}

	// Write the updated data back to the file
	tempFile := dataFile + ".tmp"
	err = os.WriteFile(tempFile, updatedData, 0644) // Changed to os.WriteFile
	if err != nil {
		log.Printf("❌ Error Writing Temp File: %v", err)
		return
	}

	// Rename the temp file to the original file
	err = os.Rename(tempFile, dataFile)
	if err != nil {
		log.Printf("❌ Error Renaming Temp File: %v", err)
		return
	}
}

func Publish(topic, message string) {
	token := client.Publish(topic, 1, false, message)
	token.Wait()
}
