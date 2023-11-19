package controllers

import (
	"SadApp/src/database"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Events(c *fiber.Ctx) error {
	return c.JSON(models.Event{})
}

func CreateEvent(c *fiber.Ctx) error {
	var event models.Event

	if err := c.BodyParser(&event); err != nil {
		return err
	}

	database.DB.Create(&event)

	return c.JSON(event)
}

func GetEvent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	event := models.Event{
		Id: uint(id),
	}

	database.DB.Find(&event)

	return c.JSON(event)
}

func UpdateEvent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	event := models.Event{
		Id: uint(id),
	}

	if err := c.BodyParser(&event); err != nil {
		return err
	}

	database.DB.Model(&event).Updates(event)

	return c.JSON(event)
}

func DeleteEvent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	event := models.Event{
		Id: uint(id),
	}

	database.DB.Delete(&event)

	return nil
}
