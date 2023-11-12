package controllers

import (
	"SadApp/src/database"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetUser 関数は、新しいユーザー詳細を取得するための関数です。
func GetUser(c *fiber.Ctx) error {
	// ユーザーモデルの新しいインスタンスを作成
	var user models.User
	// リクエストからIDパラメータを取得し、整数型に変換
	id, _ := strconv.Atoi(c.Params("id"))

	// 変換されたIDをユーザーモデルのIDに割り当て
	user.Id = uint(id)

	// 取得したidを使ってユーザーを検索
	database.DB.Find(&user)

	// ユーザーのIDと名前のみをJSON形式で返す
	return c.JSON(fiber.Map{"id": user.Id, "name": user.Name})
}
