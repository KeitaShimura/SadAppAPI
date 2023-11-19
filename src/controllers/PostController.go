package controllers

import (
	"SadApp/src/database"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Posts(c *fiber.Ctx) error {
	return c.JSON(models.Post{})
}

func CreatePost(c *fiber.Ctx) error {
	var post models.Post

	if err := c.BodyParser(&post); err != nil {
		return err
	}

	database.DB.Create(&post)

	return c.JSON(post)
}

func GetPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{
		Id: uint(id),
	}

	database.DB.Find(&post)

	return c.JSON(post)
}

func UpdatePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{
		Id: uint(id),
	}

	if err := c.BodyParser(&post); err != nil {
		return err
	}

	database.DB.Model(&post).Updates(post)

	return c.JSON(post)
}

func DeletePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{
		Id: uint(id),
	}

	database.DB.Delete(&post)

	return nil
}
