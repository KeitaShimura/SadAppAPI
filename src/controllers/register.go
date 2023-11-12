package controllers

import (
	"SadApp/src/database"         // データベースへのアクセスを提供するパッケージをインポート
	"SadApp/src/models"           // モデルを提供するパッケージをインポート
	"github.com/gofiber/fiber/v2" // Fiberフレームワークをインポート
	"golang.org/x/crypto/bcrypt"  // パスワードの暗号化に使用するパッケージをインポート
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

	// ユーザーデータをJSON形式で返す。
	return c.JSON(user)
}
