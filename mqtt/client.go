package mqtt

import (
	"fmt"
	"log"
	"os"
	"server/controllers"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var client mqtt.Client
var topics map[string]string

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
	controllers.HandleSensorData(friendlyName, string(msg.Payload()))
}

func Publish(topic, message string) {
	token := client.Publish(topic, 1, false, message)
	token.Wait()
}
