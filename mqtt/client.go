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

	payloadStr := string(msg.Payload())

	if friendlyName == "temperature" {
		var temp int
		_, err := fmt.Sscanf(payloadStr, "%d", &temp)
		if err == nil {
			temp -= 10
			payloadStr = fmt.Sprintf("%d", temp)
		} else {
			log.Printf("❌ Temperature parsing error: %v", err)
		}
	}

	sensorData := map[string]string{
		"data":    payloadStr,
		"message": fmt.Sprintf("🚪 %s Sensor Triggered: %s", friendlyName, payloadStr),
		"topic":   friendlyName,
	}

	dataJson, err := json.Marshal(sensorData)
	if err != nil {
		log.Printf("❌ Error Marshalling Sensor Data: %v", err)
		return
	}

	saveOrUpdateData(friendlyName, dataJson)
	controllers.HandleSensorData(friendlyName, payloadStr)
}

func saveOrUpdateData(topic string, data []byte) {
	mu.Lock()
	defer mu.Unlock()

	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		initialData := "{}"
		err := os.WriteFile(dataFile, []byte(initialData), 0644)
		if err != nil {
			log.Printf("❌ Error Creating Initial JSON File: %v", err)
			return
		}
	}

	fileData, err := os.ReadFile(dataFile)
	if err != nil {
		log.Printf("❌ Error Reading File: %v", err)
		return
	}

	var existingData map[string]json.RawMessage
	if len(fileData) == 0 {
		existingData = make(map[string]json.RawMessage)
	} else {
		err = json.Unmarshal(fileData, &existingData)
		if err != nil {
			log.Printf("❌ Error Unmarshalling File Data: %v", err)
			existingData = make(map[string]json.RawMessage)
		}
	}

	existingData[topic] = data

	updatedData, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		log.Printf("❌ Error Marshalling Updated Data: %v", err)
		return
	}

	tempFile := dataFile + ".tmp"
	err = os.WriteFile(tempFile, updatedData, 0644)
	if err != nil {
		log.Printf("❌ Error Writing Temp File: %v", err)
		return
	}

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
