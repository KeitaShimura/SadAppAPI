package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func PostComments(c *fiber.Ctx) error {
	postId, _ := strconv.Atoi(c.Params("id"))
	var comments []models.PostComment
	database.DB.Preload("User").Where("post_id = ?", postId).Order("created_at DESC").Find(&comments)
	return c.JSON(comments)
}

func CreatePostComment(c *fiber.Ctx) error {
	// JWTトークンからユーザーIDを取得
	userId, err := middlewares.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "認証に失敗しました。",
		})
	}

	// URLからpostIdを取得
	postId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "無効な投稿IDです。",
		})
	}

	// 新しいコメントオブジェクトを初期化
	var comment models.PostComment

	// リクエストボディをコメントオブジェクトに解析
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "不正なリクエストです。",
		})
	}

	// コメントの内容が空でないことを確認
	content := strings.TrimSpace(comment.Content)
	if len(content) == 0 || len(content) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "コメントは1文字以上500文字以下である必要があります。",
		})
	}

	// ユーザーIDと投稿IDをコメントに割り当て
	comment.UserId = userId
	comment.PostId = uint(postId)

	// コメントをデータベースに保存
	result := database.DB.Create(&comment)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "コメントの作成に失敗しました。",
		})
	}
	database.DB.Preload("User").Find(&comment, comment.Id)

	// 作成されたコメントをJSON形式で返却
	return c.JSON(comment)
}

func UpdatePostComment(c *fiber.Ctx) error {
	// URLからコメントIDを取得
	commentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "無効なコメントIDです。",
		})
	}

	// データベースからコメントを検索
	var comment models.PostComment
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

	// コメントの内容が空でないことを確認
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
			"error": "コメントの更新に失敗しました。",
		})
	}

	// 更新されたコメントデータを返却
	return c.JSON(comment)
}

func DeletePostComment(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	comment := models.PostComment{
		Id: uint(id),
	}
	database.DB.Delete(&comment)
	return nil
}
