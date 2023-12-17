package main

import (
	"log"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
	}))

	app.Post("/api/auth/signup", routes.SignUp)
	app.Post("/api/auth/login", routes.Login)

	app.Get("/api/users/:username/posts", routes.GetPosts)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	app.Post("/api/users/me/posts", routes.CreatePost)
	app.Delete("/api/users/me/posts/:postID", routes.DeletePost)

	log.Fatal(app.Listen(":8080"))
}
