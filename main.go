package main

import (
	"server/config"
	"server/mqtt"
	"server/server"
)

func main() {
	config.LoadEnv()
	mqtt.ConnectMQTT()
	server.StartServer()
}
