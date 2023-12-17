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

func matchName(name string) bool {
	return utf8.RuneCountInString(name) > 3 && utf8.RuneCountInString(name) < 30
}

func matchPassword(password string) bool {
	return utf8.RuneCountInString(password) >= 6
}

func SignUp(c *fiber.Ctx) error {
	var body struct {
		Username string
		Name     string
		Password string
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if !matchUsername(body.Username) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "BAD_USERNAME")
	}

	if !matchName(body.Name) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "BAD_NAME")
	}

	if !matchPassword(body.Password) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "BAD_PASSWORD")
	}

	var user models.User
	db.DB.First(&user, "username = ?", body.Username)

	if user.ID != 0 {
		return fiber.ErrConflict
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &models.User{
		Username: body.Username,
		Name: body.Name,
		Password: string(passwordHash),
	}
	db.DB.Create(&newUser)

	token, err := utils.GetToken(body.Username)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id":    newUser.ID,
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
