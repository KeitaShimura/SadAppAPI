package controllers

import (
	// 既存のインポート
	"SadApp/src/database"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt" // JWTを使用するためのパッケージをインポート
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

// Register 関数は、新しいユーザーを登録するための関数です。
func Register(c *fiber.Ctx) error {
	var data map[string]string // リクエストボディからデータを取得するためのマップ

	// リクエストボディの解析を試みる。エラーがあれば、それを返す。
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// パスワードとパスワード確認が一致していない場合、400ステータスコードを返す。
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "パスワードが一致しません。", // エラーメッセージをJSONで返す
		})
	}

	// パスワードを暗号化する。
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)

	// Userモデルのインスタンスを作成する。
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	// データベースにユーザーを保存する。
	database.DB.Create(&user)

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
			"message": "トークンの生成に失敗しました。",
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

	// ユーザーデータをJSON形式で返す。
	return c.JSON(user)
}
