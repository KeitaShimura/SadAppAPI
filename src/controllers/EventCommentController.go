package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func EventComments(c *fiber.Ctx) error {
	eventId, err := strconv.Atoi(c.Params("event_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	// ページ番号とページサイズを取得
	page, pageSize := getPaginationParameters(c)

	var comments []models.EventComment
	var total int64
	database.DB.Model(&models.EventComment{}).Where("event_id = ?", eventId).Count(&total)

	result := database.DB.Preload("User").
		Where("event_id = ?", eventId).
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&comments)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve comments for the event",
		})
	}

	return c.JSON(comments)
}

func CreateEventComment(c *fiber.Ctx) error {
	// JWTトークンからユーザーIDを取得
	userId, err := middlewares.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "認証に失敗しました。",
		})
	}

	// URLからeventIdを取得
	eventId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "無効なイベントIDです。",
		})
	}

	// 新しいイベントコメントオブジェクトを初期化
	var comment models.EventComment

	// リクエストボディをコメントオブジェクトに解析
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "不正なリクエストです。",
		})
	}

	// コメント内容のバリデーション
	content := strings.TrimSpace(comment.Content)
	if len(content) == 0 || len(content) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "コメントは1文字以上500文字以下である必要があります。",
		})
	}

	// ユーザーIDとイベントIDをコメントに割り当て
	comment.UserId = userId
	comment.EventId = uint(eventId)

	// コメントをデータベースに保存
	result := database.DB.Create(&comment)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "コメントを作成できませんでした。",
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
			"error": "無効なコメントIDです。",
		})
	}

	// データベースからコメントを検索
	var comment models.EventComment
	result := database.DB.Preload("User").First(&comment, commentId)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "コメントが見つかりません。",
		})
	}

	// リクエストボディから更新データを解析
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "不正なリクエストです。",
		})
	}

	// コメント内容のバリデーション
	content := strings.TrimSpace(comment.Content)
	if len(content) == 0 || len(content) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "コメントは1文字以上500文字以下である必要があります。",
		})
	}

	// コメントを更新
	result = database.DB.Save(&comment)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "コメントを更新できませんでした。",
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
