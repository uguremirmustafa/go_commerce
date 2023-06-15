package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/uguremirmustafa/go_commerce/config"
	"github.com/uguremirmustafa/go_commerce/database"
	"github.com/uguremirmustafa/go_commerce/middlewares"
	"github.com/uguremirmustafa/go_commerce/routes"
)

func init() {
	config.LoadEnvVariables()
	database.Connect()
}

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to my api")
}

func setupRoutes(app *fiber.App) {
	// welcome
	app.Get("/api", welcome)
	// auth
	app.Post("/api/signup", routes.Singup)
	app.Post("/api/login", routes.Login)
	app.Get("/api/validate", routes.ValidateAuthorization)

	// User endpoints
	app.Post("/api/users", middlewares.RequireAuth, routes.CreateUser)
	app.Get("/api/users", middlewares.RequireAuth, routes.GetUsers)
	app.Get("/api/users/:id", routes.GetUser)
	app.Put("/api/users/:id", routes.UpdateUser)
	app.Delete("/api/users/:id", routes.DeleteUser)
}

func main() {
	app := fiber.New()
	setupRoutes(app)

	port := os.Getenv("PORT")

	log.Fatal(app.Listen(":" + port))
}
