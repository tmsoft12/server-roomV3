package server

import (
	"log"
	"os"
	"server/database"
	"server/routes"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2"
)

func StartServer() {
	database.InitDatabase()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	routes.SetupRoutes(app)

	port := os.Getenv("HTTP_PORT")
	log.Println("🌍 Fiber HTTP Server started on port", port)
	log.Fatal(app.Listen(":" + port))
}
