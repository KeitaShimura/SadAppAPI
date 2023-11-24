package controllers

import (
	"SadApp/src/database"
	"SadApp/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func EventComments(c *fiber.Ctx) error {
	eventId, _ := strconv.Atoi(c.Params("event_id"))
	var comments []models.EventComment
	database.DB.Where("event_id = ?", eventId).Find(&comments)
	return c.JSON(comments)
}

func CreateEventComment(c *fiber.Ctx) error {
	comment := new(models.EventComment)
	if err := c.BodyParser(comment); err != nil {
		return err
	}

	database.DB.Create(&comment)
	return c.JSON(comment)
}

func UpdateEventComment(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	comment := models.EventComment{
		Id: uint(id),
	}
	if err := c.BodyParser(&comment); err != nil {
		return err
	}
	database.DB.Model(&comment).Updates(comment)
	return c.JSON(comment)
}

func DeleteEventComment(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	comment := models.EventComment{
		Id: uint(id),
	}
	database.DB.Delete(&comment)
	return nil
}
