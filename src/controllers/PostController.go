package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Posts(c *fiber.Ctx) error {
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

	// Fetch posts from the current user and the users who follow them
	var posts []models.Post
	result := database.DB.Where("user_id IN ?", followerIds).
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&posts)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve posts",
		})
	}

	// Return the list of posts as JSON
	return c.JSON(posts)
}

func UserPosts(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// ページ番号とページサイズを取得
	page, pageSize := getPaginationParameters(c)

	var posts []models.Post
	var total int64
	database.DB.Model(&models.Post{}).Where("user_id = ?", userID).Count(&total)

	result := database.DB.Where("user_id = ?", userID).
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&posts)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve posts for the user",
		})
	}

	return c.JSON(posts)
}

func GetPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{
		Id: uint(id),
	}

	database.DB.Preload("User").Preload("PostComment").Find(&post)

	return c.JSON(post)
}

func CreatePost(c *fiber.Ctx) error {
	// 最初にJWTトークンからユーザーIDを取得します
	userId, err := middlewares.GetUserId(c)
	if err != nil {
		// ユーザーIDを取得できない場合はエラーを返します
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "認証に失敗しました。",
		})
	}

	// 新しいPost構造体を初期化します
	var post models.Post

	// リクエストボディをPost構造体に解析します
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "不正なリクエストです。",
		})
	}

	content := strings.TrimSpace(post.Content)
	if len(content) == 0 || len(content) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "コメントは1文字以上500文字以下である必要があります。",
		})
	}

	// 取得したユーザーIDをPostに割り当てます
	post.UserId = userId // PostモデルにUserIdフィールドがあると仮定しています

	// データベースにPostを作成します
	result := database.DB.Create(&post)
	if result.Error != nil {
		// 作成時のエラーを処理します
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "投稿の作成に失敗しました。",
		})
	}

	// 作成された投稿のUserデータを読み込みます
	database.DB.Preload("User").Find(&post, post.Id)

	// 作成された投稿をJSON形式で返します
	return c.JSON(post)
}

func UpdatePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{
		Id: uint(id),
	}

	content := strings.TrimSpace(post.Content)
	if len(content) == 0 || len(content) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "コメントは1文字以上500文字以下である必要があります。",
		})
	}

	if err := c.BodyParser(&post); err != nil {
		return err
	}

	// Postを更新します
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

func UserLikedPosts(c *fiber.Ctx) error {
	// Retrieve the user ID (adjust this part based on how you manage user identification)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get pagination parameters (reuse your existing function or define one)
	page, pageSize := getPaginationParameters(c)

	// Find IDs of posts liked by the user
	var postLikes []models.PostLike
	database.DB.Where("user_id = ?", userID).Find(&postLikes)

	// Extract post IDs
	var postIds []uint
	for _, postLike := range postLikes {
		postIds = append(postIds, postLike.PostId)
	}

	// Fetch the posts based on the post IDs
	var posts []models.Post
	result := database.DB.Where("id IN ?", postIds).
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&posts)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve liked posts",
		})
	}

	// Return the list of liked posts as JSON
	return c.JSON(posts)
}

// Helper function to get pagination parameters
func getPaginationParameters(c *fiber.Ctx) (int, int) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "100"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}
