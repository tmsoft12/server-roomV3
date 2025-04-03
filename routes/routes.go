package routes

import (
	"server/controllers"
	"server/mqtt"
	wscontroller "server/wsController"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "MQTT Service is running 🚀"})
	})

	app.Post("/send", func(c *fiber.Ctx) error {
		type Request struct {
			Topic   string `json:"topic"`
			Message string `json:"message"`
		}
		req := new(Request)
		if err := c.BodyParser(req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		if req.Topic == "" || req.Message == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Missing topic or message"})
		}
		mqtt.Publish(req.Topic, req.Message)
		return c.JSON(fiber.Map{"message": "Message sent successfully"})
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		wscontroller.Register(c)
	}))
	app.Get("/logg", controllers.HandleGetSensorsWithPagination)
	app.Get("/search", controllers.GetSensorDataByDateHandler)
}
