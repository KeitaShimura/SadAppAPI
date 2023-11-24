package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// イベントへの参加
func JoinEvent(c *fiber.Ctx) error {
	// ログインユーザーのIDを取得
	authUserId, _ := middlewares.GetUserId(c)

	var participant models.EventParticipant

	// パラメータからevent_idを取得
	eventId, _ := strconv.Atoi(c.Params("id"))
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
	eventId, _ := strconv.Atoi(c.Params("id"))

	// 該当するイベント参加情報を検索して削除
	participant := models.EventParticipant{}
	database.DB.Where("event_id = ? AND user_id = ?", eventId, authUserId).Delete(&participant)

	return c.SendStatus(fiber.StatusOK)
}

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
