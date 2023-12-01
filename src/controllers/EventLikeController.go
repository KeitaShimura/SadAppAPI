package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// 投稿に「いいね」を追加する
func LikeEvent(c *fiber.Ctx) error {
	userId, _ := middlewares.GetUserId(c)

	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	like := models.EventLike{
		UserId:  userId,
		EventId: uint(eventId),
	}

	database.DB.Create(&like)

	return c.JSON(like)
}

// 投稿への「いいね」を削除する
func UnlikeEvent(c *fiber.Ctx) error {
	userId, _ := middlewares.GetUserId(c)

	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	database.DB.Where("user_id = ? AND event_id = ?", userId, eventId).Delete(&models.EventLike{})

	return c.SendStatus(fiber.StatusOK)
}

func CheckIfEventLiked(c *fiber.Ctx) error {
	userId, _ := middlewares.GetUserId(c)
	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	var like models.EventLike
	result := database.DB.Where("user_id = ? AND event_id = ?", userId, eventId).First(&like)

	if result.Error != nil {
		return c.JSON(false)
	}

	return c.JSON(true)
}

func GetLikesForEvent(c *fiber.Ctx) error {
	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	var likes []models.EventLike
	database.DB.Where("event_id = ?", eventId).Order("created_at DESC").Find(&likes)

	return c.JSON(likes)
}
