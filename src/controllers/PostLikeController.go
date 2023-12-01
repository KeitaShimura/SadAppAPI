package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// 投稿に「いいね」を追加する
func LikePost(c *fiber.Ctx) error {
	userId, _ := middlewares.GetUserId(c)

	postId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	like := models.PostLike{
		UserId: userId,
		PostId: uint(postId),
	}

	database.DB.Create(&like)

	return c.JSON(like)
}

// 投稿への「いいね」を削除する
func UnlikePost(c *fiber.Ctx) error {
	userId, _ := middlewares.GetUserId(c)

	postId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	database.DB.Where("user_id = ? AND post_id = ?", userId, postId).Delete(&models.PostLike{})

	return c.SendStatus(fiber.StatusOK)
}

func CheckIfPostLiked(c *fiber.Ctx) error {
	userId, _ := middlewares.GetUserId(c)
	postId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	var like models.PostLike
	result := database.DB.Where("user_id = ? AND post_id = ?", userId, postId).First(&like)

	if result.Error != nil {
		return c.JSON(false)
	}

	return c.JSON(true)
}

func GetLikesForPost(c *fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	var likes []models.PostLike
	database.DB.Where("post_id = ?", postId).Order("created_at DESC").Find(&likes)

	return c.JSON(likes)
}
