package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Events(c *fiber.Ctx) error {
	currentUserId, err := middlewares.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving user ID",
		})
	}

	// Fetch the IDs of users that are following the current user
	var followers []models.Follow
	database.DB.Where("following_id = ?", currentUserId).Find(&followers)

	// Extract user IDs from the followers
	var followerIds []uint
	for _, follower := range followers {
		followerIds = append(followerIds, follower.FollowerId)
	}

	// Include the current user's ID in the list
	followerIds = append(followerIds, currentUserId)

	// Get pagination parameters
	page, pageSize := getPaginationParameters(c)

	// Fetch events from the current user and the users who follow them
	var events []models.Event
	result := database.DB.Where("user_id IN ?", followerIds).
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&events)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve events",
		})
	}

	// Return the list of events as JSON
	return c.JSON(events)
}

func UserEvents(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var events []models.Event
	// ページ番号を取得
	page := 1
	pageSize := 100

	// クエリから 'page' を取得
	if p, err := strconv.Atoi(c.Query("page", "1")); err == nil && p > 0 {
		page = p
	}

	// クエリから 'pageSize' を取得
	if ps, err := strconv.Atoi(c.Query("pageSize", "100")); err == nil && ps > 0 {
		pageSize = ps
	}

	var total int64
	database.DB.Model(&models.Event{}).Where("user_id = ?", userID).Count(&total)

	result := database.DB.Where("user_id = ?", userID).
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&events)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve events for the user",
		})
	}

	return c.JSON(events)
}

func GetEvent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	event := models.Event{
		Id: uint(id),
	}

	database.DB.Preload("User").Preload("EventComment").Find(&event)

	return c.JSON(event)
}

func CreateEvent(c *fiber.Ctx) error {
	// まず、JWTトークンからユーザーIDを取得
	userId, err := middlewares.GetUserId(c)
	if err != nil {
		// ユーザーIDを取得できない場合、エラーを返す
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "認証に失敗しました。",
		})
	}

	// 新しいイベント構造体を初期化
	var event models.Event

	// リクエストボディをイベント構造体に解析
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "不正なリクエストです。",
		})
	}

	// 画像ファイルの処理
	file, err := c.FormFile("image")
	if err == nil {
		// 安全なファイル名の生成
		fileName := filepath.Base(file.Filename)
		safeFileName := fmt.Sprintf("%d-%s", time.Now().Unix(), fileName)

		// 保存先パスの生成
		imagePath := filepath.Join("src/uploads", safeFileName)

		// ディレクトリの存在確認と作成
		if _, err := os.Stat("src/uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("src/uploads", 0755); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "画像保存用のディレクトリの作成に失敗しました。",
				})
			}
		}

		// 画像の保存
		if err := c.SaveFile(file, imagePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "画像の保存に失敗しました。",
			})
		}

		// 画像のURLをイベントに割り当て
		event.Image = "/" + imagePath
	}

	// イベントタイトルのバリデーション
	if len(event.Title) == 0 || len(event.Title) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "イベントタイトルは1文字以上100文字以下である必要があります。",
		})
	}

	// イベント内容のバリデーション
	if len(event.Content) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "イベント内容は500文字以下である必要があります。",
		})
	}

	// イベントURLのバリデーション
	content := strings.TrimSpace(event.Content)
	if len(content) == 0 || len(content) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "コメントは1文字以上500文字以下である必要があります。",
		})
	}

	// Assign the retrieved user ID to the event
	event.UserId = userId // Assuming your event model has a UserId field

	// データベースにイベントを保存
	result := database.DB.Create(&event)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "イベントを作成できませんでした。",
		})
	}

	// 作成されたイベントのUserデータを読み込み
	database.DB.Preload("User").Find(&event, event.Id)

	// 作成されたイベントをJSON形式で返す
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

	// 更新データのバリデーション
	if len(event.Title) == 0 || len(event.Title) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "イベントタイトルは1文字以上100文字以下である必要があります。",
		})
	}

	content := strings.TrimSpace(event.Content)
	if len(content) == 0 || len(content) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "コメントは1文字以上500文字以下である必要があります。",
		})
	}

	if len(event.Event_URL) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "イベントURLは255文字以下である必要があります。",
		})
	}

	database.DB.Model(&event).Updates(event)

	return c.JSON(event)
}

func UserLikedEvents(c *fiber.Ctx) error {
	// Retrieve the user ID (adjust this part based on how you manage user identification)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get pagination parameters (reuse your existing function or define one)
	page, pageSize := getPaginationParameters(c)

	// Find IDs of events liked by the user
	var eventLikes []models.EventLike
	database.DB.Where("user_id = ?", userID).Find(&eventLikes)

	// Extract event IDs
	var eventIds []uint
	for _, eventLike := range eventLikes {
		eventIds = append(eventIds, eventLike.EventId)
	}

	// Fetch the events based on the event IDs
	var events []models.Event
	result := database.DB.Where("id IN ?", eventIds).
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&events)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve liked events",
		})
	}

	// Return the list of liked events as JSON
	return c.JSON(events)
}

func UserParticipatedEvents(c *fiber.Ctx) error {
	// Retrieve the user ID (adjust this part based on how you manage user identification)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get pagination parameters (reuse your existing function or define one)
	page, pageSize := getPaginationParameters(c)

	// Find IDs of events liked by the user
	var eventParticipants []models.EventParticipant
	database.DB.Where("user_id = ?", userID).Find(&eventParticipants)

	// Extract event IDs
	var eventIds []uint
	for _, eventParticipant := range eventParticipants {
		eventIds = append(eventIds, eventParticipant.EventId)
	}

	// Fetch the events based on the event IDs
	var events []models.Event
	result := database.DB.Where("id IN ?", eventIds).
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&events)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve participated events",
		})
	}

	// Return the list of participated events as JSON
	return c.JSON(events)
}

func DeleteEvent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	event := models.Event{
		Id: uint(id),
	}

	database.DB.Delete(&event)

	return nil
}
