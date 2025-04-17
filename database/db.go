package database

import (
	"database/sql"
	"fmt"
	"log"
	"server/models"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite", "./sensor_data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sensor_data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		topic TEXT,
		data TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("✅ SQLite Database Initialized")
}

func InsertSensorData(topic, data string) error {
	if data == "open" {
		data = " Açyldy"
	} else if data == "closed" {
		data = " Ýapyldy"
	}
	if "21" <= data && topic == "temperature" {
		topic = "Gyzgynlyk"
		data = fmt.Sprintf("%s C", data)
		add(topic, data)
	}
	if "31" <= data && topic == "humidity" {
		topic = "Çyglylyk"
		data = fmt.Sprintf("%s %s", data, string("%"))
		add(topic, data)
	}

	if topic == "main" {
		topic = "Easasy Gapy"
		add(topic, data)
	}

	switch topic {
	case "door1":
		topic = "Gapy 1"
		add(topic, data)

	case "door2":
		topic = "Gapy 2"
		add(topic, data)

	case "door3":
		topic = "Gapy 3"
		add(topic, data)

	case "door4":
		topic = "Gapy 4"
		add(topic, data)

	case "door5":
		topic = "Gapy 5"
		add(topic, data)

	case "door6":
		topic = "Gapy 6"
		add(topic, data)

	case "door7":
		topic = "Gapy 7"
		add(topic, data)

	}

	return nil
}
func add(topic string, data string) error {

	insertSQL := `INSERT INTO sensor_data (topic, data) VALUES (?, ?)`
	_, err := DB.Exec(insertSQL, topic, data)
	if err != nil {
		return fmt.Errorf("failed to insert sensor data: %v", err)
	}
	return nil
}

type Sensor struct {
	ID        int    `json:"id"`
	Topic     string `json:"topic"`
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

func GetSensorDataWithPagination(page, pageSize int) ([]Sensor, int64, error) {
	offset := (page - 1) * pageSize
	var totalRecords int64

	err := DB.QueryRow("SELECT COUNT(*) FROM sensor_data").Scan(&totalRecords)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch total records: %v", err)
	}

	rows, err := DB.Query("SELECT id, topic, data, timestamp FROM sensor_data ORDER BY id DESC LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch sensor data: %v", err)
	}
	defer rows.Close()

	var sensors []Sensor
	for rows.Next() {
		var s Sensor
		if err := rows.Scan(&s.ID, &s.Topic, &s.Data, &s.Timestamp); err != nil {
			return nil, 0, fmt.Errorf("failed to scan sensor data: %v", err)
		}
		sensors = append(sensors, s)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during row iteration: %v", err)
	}

	return sensors, totalRecords, nil
}

func CloseDatabase() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		} else {
			log.Println("✅ Database connection closed")
		}
	}
}
func GetSensorDataByDate(startDate, endDate string, page, pageSize int) ([]models.Sensor, int64, int, bool, bool, error) {
	if page < 1 || pageSize < 1 {
		return nil, 0, 0, false, false, fmt.Errorf("invalid pagination parameters: page=%d, pageSize=%d", page, pageSize)
	}
	if pageSize > 1000 {
		return nil, 0, 0, false, false, fmt.Errorf("pageSize cannot exceed 1000")
	}

	const dateFormat = "2006-01-02"
	if _, err := time.Parse(dateFormat, startDate); err != nil {
		return nil, 0, 0, false, false, fmt.Errorf("invalid start date format: %v", err)
	}
	if _, err := time.Parse(dateFormat, endDate); err != nil {
		return nil, 0, 0, false, false, fmt.Errorf("invalid end date format: %v", err)
	}

	offset := (page - 1) * pageSize
	var totalRecords int64

	countQuery := "SELECT COUNT(*) FROM sensor_data WHERE DATE(timestamp) BETWEEN ? AND ?"
	err := DB.QueryRow(countQuery, startDate, endDate).Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, false, false, fmt.Errorf("failed to fetch total records: %v", err)
	}

	totalPages := int(totalRecords / int64(pageSize))
	if totalRecords%int64(pageSize) != 0 {
		totalPages++
	}

	hasNext := page < totalPages
	hasPrev := page > 1

	query := `SELECT id, topic, data, timestamp 
              FROM sensor_data 
              WHERE DATE(timestamp) BETWEEN ? AND ? 
              ORDER BY id DESC 
              LIMIT ? OFFSET ?`

	rows, err := DB.Query(query, startDate, endDate, pageSize, offset)
	if err != nil {
		return nil, 0, 0, false, false, fmt.Errorf("failed to fetch sensor data: %v", err)
	}
	defer rows.Close()

	var sensors []models.Sensor
	for rows.Next() {
		var s models.Sensor
		if err := rows.Scan(&s.ID, &s.Topic, &s.Data, &s.Timestamp); err != nil {
			return nil, 0, 0, false, false, fmt.Errorf("failed to scan sensor data: %v", err)
		}
		sensors = append(sensors, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, false, false, fmt.Errorf("error during row iteration: %v", err)
	}

	return sensors, totalRecords, totalPages, hasNext, hasPrev, nil
}
