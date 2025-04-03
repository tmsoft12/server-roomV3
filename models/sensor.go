package models

type Sensor struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Topic     string `json:"topic"`
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}
