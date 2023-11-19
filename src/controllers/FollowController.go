package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Follow(c *fiber.Ctx) error {
	// ログインユーザーのIDを取得
	authUserId, _ := middlewares.GetUserId(c)

	var follow models.Follow

	if err := c.BodyParser(&follow); err != nil {
		return err
	}

	// ログインユーザーのIDをfollowing_idに設定
	follow.FollowerId = authUserId

	// フォロー情報をデータベースに保存
	database.DB.Create(&follow)

	return c.JSON(follow)
}

func UnFollow(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	follow := models.Follow{
		Id: uint(id),
	}

	database.DB.Delete(&follow)

	return nil
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
