package main

import (
	"log"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
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

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

    log.Fatal(app.Listen(":8080"))
}
