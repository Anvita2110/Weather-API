package main

import (
	"os"
	"os/signal"
	"syscall"
	"weather-app/routers"

	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3001"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	routers.SetupRoutes(app)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		routers.CloseRedisClient()

		app.Shutdown()
	}()

	app.Listen(":3000")
}
