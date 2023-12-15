package routes

import (
	"regexp"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/real-firesnap/firesnap-api/db"
	"github.com/real-firesnap/firesnap-api/models"
	"github.com/real-firesnap/firesnap-api/utils"
	"golang.org/x/crypto/bcrypt"
)

var matchUsername = regexp.MustCompile(`^[a-zA-Z0-9\-_]{3,15}$`).MatchString

func matchPassword(password string) bool {
	return utf8.RuneCountInString(password) >= 6
}

func SignUp(c *fiber.Ctx) error {
	var body struct {
		Username string
		Password string
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if !matchUsername(body.Username) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "BAD_USERNAME")
	}

	if !matchPassword(body.Password) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "BAD_PASSWORD")
	}

	var user models.User
	db.DB.First(&user, "username = ?", body.Username)

	if user.ID != 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "USERNAME_TAKEN")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	db.DB.Create(&models.User{
		Username: body.Username,
		Password: string(passwordHash),
	})

	token, err := utils.GetToken(body.Username)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func Login(c *fiber.Ctx) error {
	var body struct {
		Username string
		Password string
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if !matchUsername(body.Username) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "BAD_USERNAME")
	}

	if !matchPassword(body.Password) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "BAD_PASSWORD")
	}

	var user models.User
	db.DB.First(&user, "username = ?", body.Username)

	if user.ID == 0 {
		return fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return fiber.ErrUnauthorized
	}

	token, err := utils.GetToken(body.Username)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
