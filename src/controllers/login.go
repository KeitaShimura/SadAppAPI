package controllers

import (
	"SadApp/src/database"
	"SadApp/src/models"           // モデルを提供するパッケージをインポート
	"github.com/gofiber/fiber/v2" // Fiberフレームワークをインポート
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

// Login 関数は、ユーザーのログイン処理を行う関数です。
func Login(c *fiber.Ctx) error {
	var data map[string]string // リクエストボディからデータを取得するためのマップ

	// リクエストボディの解析を試みる。エラーがあれば、それを返す。
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User // ユーザーモデルのインスタンス

	// データベースからメールアドレスを使ってユーザー情報を検索する。
	database.DB.Where("email = ?", data["email"]).First(&user)

	// ユーザーが見つからなかった場合、エラーを返す。
	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "認証に失敗しました。", // 認証失敗のメッセージをJSONで返す
		})
	}

	// パスワードが一致するかチェックする。一致しない場合はエラーを返す。
	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "認証に失敗しました。", // 認証失敗のメッセージをJSONで返す
		})
	}

	// JWTトークンの生成
	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	// JWTトークンを生成し、エラーがあれば処理を返す。
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "認証に失敗しました。",
		})
	}

	// クッキーにJWTトークンを設定する。
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	// 認証が成功したユーザーデータをJSON形式で返す。
	return c.JSON(user)
}
