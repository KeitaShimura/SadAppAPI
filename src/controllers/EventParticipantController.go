package controllers

import (
	"SadApp/src/database"
	"SadApp/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetEventParticipants(c *fiber.Ctx) error {
	// イベントIDをURLパラメータから取得
	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "無効なイベントID"})
	}

	// EventParticipant モデルを使用して、特定のイベントIDに対する参加者を取得
	var participants []models.EventParticipant
	database.DB.Where("event_id = ?", eventId).Find(&participants)

	// 取得した参加者のリストをJSONとして返す
	return c.JSON(participants)
}
