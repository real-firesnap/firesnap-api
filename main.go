package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/real-firesnap/firesnap-api/db"
	"github.com/real-firesnap/firesnap-api/routes"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDatabase()
	db.MigrateDatabase()
}

func main() {
	app := fiber.New()

    app.Post("/api/auth/signup", routes.SignUp)
	app.Post("/api/auth/login", routes.Login)

    log.Fatal(app.Listen(":8080"))
}
