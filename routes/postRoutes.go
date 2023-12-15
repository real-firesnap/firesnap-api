package routes

import (
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/real-firesnap/firesnap-api/db"
	"github.com/real-firesnap/firesnap-api/models"
	"github.com/real-firesnap/firesnap-api/utils"
)

func CreatePost(c *fiber.Ctx) error {
	var body struct {
		Content string
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if body.Content == "" || utf8.RuneCountInString(body.Content) > 1000 {
		return fiber.NewError(fiber.ErrUnprocessableEntity.Code, "BAD_CONTENT_LENGTH")
	}

	user, err := utils.AuthUser(c)
	if err != nil {
		return err
	}

	post := models.Post{
		User: user,
		Content: body.Content,
	}

	db.DB.Create(&post)

	return c.JSON(fiber.Map{
		"postID": post.ID,
	})
}

func GetPosts(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	db.DB.First(&user, "username = ?", username)

	if user.ID == 0 {
		return fiber.ErrNotFound
	}

	var posts []models.Post
	db.DB.Limit(30).Order("created_at DESC").Find(&posts, "user_id = ?", user.ID)

	result := make([]fiber.Map, 0)

	for _, post := range posts {
		result = append(result, fiber.Map{
			"id": post.ID,
			"createdAt": post.CreatedAt.Unix(),
			"content": post.Content,
		})
	}

	return c.JSON(result)
}

func DeletePost(c *fiber.Ctx) error {
	postID := c.Params("postID")

	user, err := utils.AuthUser(c)
	if err != nil {
		return err
	}

	var post models.Post
	db.DB.First(&post, "id = ? AND user_id = ?", postID, user.ID)

	if post.ID == 0 {
		return fiber.ErrNotFound
	}

	db.DB.Delete(&post, post.ID)

	return nil
}
