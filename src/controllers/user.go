package controllers

import (
	"SadApp/src/database"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"strconv"
)

func GetAuthUser(c *fiber.Ctx) error {
	// ユーザーのブラウザから"jwt"という名前のCookieを取得
	cookie := c.Cookies("jwt")

	// 取得したCookieを使用してJWTを解析
	// ここでは"secret"をキーとして使用している
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	// JWTが無効であるか、解析中にエラーが発生した場合、認証エラーを返す
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "認証されていません。", // "Unauthenticated"のメッセージをユーザーに表示
		})
	}

	// 解析したトークンからユーザーのクレーム情報を取得
	payload := token.Claims.(*jwt.StandardClaims)

	// ユーザーモデルの新しいインスタンスを作成
	var user models.User

	// データベースからユーザーIDに基づいてユーザー情報を取得
	// payload.SubjectにはユーザーIDが含まれている
	database.DB.Where("id = ?", payload.Subject).First(&user)

	// ユーザー情報をJSON形式で返す
	return c.JSON(user)
}

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
