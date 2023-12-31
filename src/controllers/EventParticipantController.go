package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// イベントへの参加
func ParticipationEvent(c *fiber.Ctx) error {
	// ログインユーザーのIDを取得
	authUserId, _ := middlewares.GetUserId(c)

	var participant models.EventParticipant

	// パラメータからevent_idを取得
	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}
	participant.EventId = uint(eventId)
	participant.UserId = authUserId

	// イベント参加情報をデータベースに保存
	database.DB.Create(&participant)

	return c.JSON(participant)
}

// イベント参加の解除
func LeaveEvent(c *fiber.Ctx) error {
	// ログインユーザーのIDを取得
	authUserId, _ := middlewares.GetUserId(c)

	// パラメータからevent_idを取得
	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}
	// 該当するイベント参加情報を検索して削除
	participant := models.EventParticipant{}
	database.DB.Where("event_id = ? AND user_id = ?", eventId, authUserId).Delete(&participant)

	return c.SendStatus(fiber.StatusOK)
}

func CheckIfEventParticipated(c *fiber.Ctx) error {
	userId, _ := middlewares.GetUserId(c)
	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	var participant models.EventParticipant
	result := database.DB.Where("user_id = ? AND event_id = ?", userId, eventId).First(&participant)

	if result.Error != nil {
		return c.JSON(false)
	}

	return c.JSON(true)
}

func GetEventParticipants(c *fiber.Ctx) error {
	// イベントIDをURLパラメータから取得
	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "無効なイベントID"})
	}

	// EventParticipant モデルを使用して、特定のイベントIDに対する参加者を取得
	var participants []models.EventParticipant
	database.DB.Where("event_id = ?", eventId).Order("created_at DESC").Find(&participants)

	// 取得した参加者のリストをJSONとして返す
	return c.JSON(participants)
}
