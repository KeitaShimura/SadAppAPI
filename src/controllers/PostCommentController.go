package controllers

import (
	"SadApp/src/database"
	"SadApp/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Comments(c *fiber.Ctx) error {
	postId, _ := strconv.Atoi(c.Params("post_id"))
	var comments []models.PostComment
	database.DB.Where("post_id = ?", postId).Find(&comments)
	return c.JSON(comments)
}

func CreateComment(c *fiber.Ctx) error {
	comment := new(models.PostComment)
	if err := c.BodyParser(comment); err != nil {
		return err
	}

	database.DB.Create(&comment)
	return c.JSON(comment)
}

func UpdateComment(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	comment := models.PostComment{
		Id: uint(id),
	}
	if err := c.BodyParser(&comment); err != nil {
		return err
	}
	database.DB.Model(&comment).Updates(comment)
	return c.JSON(comment)
}

func DeleteComment(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	comment := models.PostComment{
		Id: uint(id),
	}
	database.DB.Delete(&comment)
	return nil
}
