package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
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
	// JWTトークンからユーザーIDを取得
	userId, err := middlewares.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// URLからeventIdを取得
	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	// 新しいイベントコメントオブジェクトを初期化
	var comment models.EventComment

	// リクエストボディをコメントオブジェクトに解析
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad request",
		})
	}

	// ユーザーIDとイベントIDをコメントに割り当て
	comment.UserId = userId
	comment.EventId = uint(eventId)

	// コメントをデータベースに保存
	result := database.DB.Create(&comment)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create comment",
		})
	}

	database.DB.Preload("User").Find(&comment, comment.Id)

	// 作成されたコメントをJSON形式で返却
	return c.JSON(comment)
}

func UpdateEventComment(c *fiber.Ctx) error {
	// URLからコメントIDを取得
	commentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	// データベースからコメントを検索
	var comment models.EventComment
	result := database.DB.Preload("User").First(&comment, commentId)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Comment not found",
		})
	}

	// リクエストボディから更新データを解析
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad request",
		})
	}

	// コメントを更新
	result = database.DB.Save(&comment)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update comment",
		})
	}

	// 更新されたコメントデータを返却
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
