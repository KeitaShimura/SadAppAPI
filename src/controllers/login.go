package controllers

import (
	"SadApp/src/database"
	"SadApp/src/models" // モデルを提供するパッケージをインポート

	"github.com/gofiber/fiber/v2" // Fiberフレームワークをインポート
	"golang.org/x/crypto/bcrypt"  // パスワードの暗号化に使用するパッケージをインポート
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
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "認証に失敗しました。", // 認証失敗のメッセージをJSONで返す
		})
	}

	// 認証が成功したユーザーデータをJSON形式で返す。
	return c.JSON(user)
}
