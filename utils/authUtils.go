package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/real-firesnap/firesnap-api/db"
	"github.com/real-firesnap/firesnap-api/models"
)

func AuthUser(c *fiber.Ctx) (models.User, error) {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	username := claims["sub"].(string)

	var user models.User
	db.DB.First(&user, "username = ?", username)

	if user.ID == 0 {
		return models.User{}, fiber.ErrUnauthorized
	}

	return user, nil
}
