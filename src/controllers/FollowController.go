package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Follow(c *fiber.Ctx) error {
	// ログインユーザーのIDを取得
	authUserId, _ := middlewares.GetUserId(c)

	var follow models.Follow

	if err := c.BodyParser(&follow); err != nil {
		return err
	}

	// ログインユーザーのIDをfollowing_idに設定
	follow.FollowingId = authUserId

	// パラメータからfollower_idを取得
	followerId, _ := strconv.Atoi(c.Params("id"))
	follow.FollowerId = uint(followerId)

	// フォロー情報をデータベースに保存
	database.DB.Create(&follow)

	return c.JSON(follow)
}

func UnFollow(c *fiber.Ctx) error {
	// ログインユーザーのIDを取得
	authUserId, _ := middlewares.GetUserId(c)

	// パラメータからfollower_idを取得
	followerId, _ := strconv.Atoi(c.Params("id"))

	// 該当するフォロー関係を検索して削除
	follow := models.Follow{}
	database.DB.Where("following_id = ? AND follower_id = ?", authUserId, followerId).Delete(&follow)

	return c.SendStatus(fiber.StatusOK)
}

func CheckIfFollowing(c *fiber.Ctx) error {
	authUserId, _ := middlewares.GetUserId(c)
	targetUserId, _ := strconv.Atoi(c.Params("id"))

	var follow models.Follow
	result := database.DB.Where("following_id = ? AND follower_id = ?", authUserId, targetUserId).First(&follow)

	if result.Error != nil {
		return c.JSON(false)
	}

	return c.JSON(true)
}

func GetFollowers(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("userID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "無効なユーザーID"})
	}

	var followers []models.User
	database.DB.Where("following_id = ?", userID).Find(&followers)

	return c.JSON(followers)
}

func GetFollowings(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("userID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "無効なユーザーID"})
	}

	var following []models.User
	database.DB.Where("follower_id = ?", userID).Find(&following)

	return c.JSON(following)
}
