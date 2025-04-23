package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"server/database"
	"server/wsController"

	"github.com/gofiber/fiber/v2"
)

func SendSMS(phoneNumber, message string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`echo "%s" | gammu --sendsms TEXT %s`, message, phoneNumber))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}
	fmt.Println("SMS Dont send:", string(output))
	return nil
}

func HandleSensorData(topic string, data string) {
	var message string
	var jsonData []byte
	phoneNumber := os.Getenv("PHONE_NUMBER")
	if phoneNumber == "" {
		phoneNumber = "+99362805208"
	}

	switch topic {
	case "temperature":
		message = fmt.Sprintf("🌡️ Temperature Data Received: %s", data)

	case "humidity":
		message = fmt.Sprintf("💧 Humidity Data Received: %s", data)

	case "main":
		message = fmt.Sprintf("🚪 Main Door Sensor Triggered: %s", data)
		if data == "open" {
			go func() {
				msg := "Server Esasy gapysy açyldy. Giriş kontrolyny we howpsuzlyk çärelerini gözden geçiriň."
				if err := SendSMS(phoneNumber, msg); err != nil {
					log.Println("SMS err:", err)
				}
			}()
		}

	case "door1":
		message = fmt.Sprintf("🚪 Door 1 Sensor Triggered: %s", data)
		if data == "open" {
			go func() {
				msg := "Server 1-nji gapysy açyldy. Giriş kontrolyny we howpsuzlyk çärelerini gözden geçiriň."
				if err := SendSMS(phoneNumber, msg); err != nil {
					log.Println("SMS err:", err)
				}
			}()
		}

	case "door2":
		message = fmt.Sprintf("🚪 Door 2 Sensor Triggered: %s", data)
		if data == "open" {
			go func() {
				msg := "Server 2-nji gapysy açyldy. Giriş kontrolyny we howpsuzlyk çärelerini gözden geçiriň."
				if err := SendSMS(phoneNumber, msg); err != nil {
					log.Println("SMS err:", err)
				}
			}()
		}

	case "door3":
		message = fmt.Sprintf("🚪 Door 3 Sensor Triggered: %s", data)
		if data == "open" {
			go func() {
				msg := "Server 3-nji gapysy açyldy. Giriş kontrolyny we howpsuzlyk çärelerini gözden geçiriň."
				if err := SendSMS(phoneNumber, msg); err != nil {
					log.Println("SMS err:", err)
				}
			}()
		}

	case "door4":
		message = fmt.Sprintf("🚪 Door 4 Sensor Triggered: %s", data)
		if data == "open" {
			go func() {
				msg := "Server 4-nji gapysy açyldy. Giriş kontrolyny we howpsuzlyk çärelerini gözden geçiriň."
				if err := SendSMS(phoneNumber, msg); err != nil {
					log.Println("SMS err:", err)
				}
			}()
		}

	case "door5":
		message = fmt.Sprintf("🚪 Door 5 Sensor Triggered: %s", data)
		if data == "open" {
			go func() {
				msg := "Server 5-nji gapysy açyldy. Giriş kontrolyny we howpsuzlyk çärelerini gözden geçiriň."
				if err := SendSMS(phoneNumber, msg); err != nil {
					log.Println("SMS err:", err)
				}
			}()
		}

	case "door6":
		message = fmt.Sprintf("🚪 Door 6 Sensor Triggered: %s", data)
		if data == "open" {
			go func() {
				msg := "Server 6-njy gapysy açyldy. Giriş kontrolyny we howpsuzlyk çärelerini gözden geçiriň."
				if err := SendSMS(phoneNumber, msg); err != nil {
					log.Println("SMS err:", err)
				}
			}()
		}

	case "door7":
		message = fmt.Sprintf("🚪 Door 7 Sensor Triggered: %s", data)
		if data == "open" {
			go func() {
				msg := "Server 7-nji gapysy açyldy. Giriş kontrolyny we howpsuzlyk çärelerini gözden geçiriň."
				if err := SendSMS(phoneNumber, msg); err != nil {
					log.Println("SMS err:", err)
				}
			}()
		}

	default:
		message = fmt.Sprintf("⚠️ Unknown Topic: %s", topic)
		log.Println(message)
	}

	payload := map[string]interface{}{
		"topic":   topic,
		"data":    data,
		"message": message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("JSON Marshalling Error:", err)
		return
	}

	wsController.Broadcast(string(jsonData))

	err = database.InsertSensorData(topic, data)
	if err != nil {
		log.Println("Error inserting sensor data into SQLite:", err)
	} else {
		log.Println("Sensor data inserted into SQLite")
	}
}

func GetSensorDataByDateHandler(c *fiber.Ctx) error {
	start := c.Query("start")
	end := c.Query("end")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	if start == "" || end == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "start and end dates are required (YYYY-MM-DD)",
		})
	}

	sensors, total, totalPages, hasNext, hasPrev, err := database.GetSensorDataByDate(start, end, page, pageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to retrieve sensor data: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": totalPages,
		"hasNext":    hasNext,
		"hasPrev":    hasPrev,
		"data":       sensors,
	})
}
