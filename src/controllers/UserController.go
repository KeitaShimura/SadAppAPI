package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetUser(c *fiber.Ctx) error {
	// ユーザーモデルの新しいインスタンスを作成
	var user models.User
	// リクエストからIDパラメータを取得し、整数型に変換
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		// IDパラメータの変換エラー処理
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "無効なIDフォーマット"})
	}

	// 変換されたIDをユーザーモデルのIDに割り当て
	user.Id = uint(id)

	// 取得したidを使ってユーザーを検索
	result := database.DB.Find(&user)

	// ユーザーが見つからなかった場合の処理
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "ユーザーが見つかりません",
		})
	}

	// ユーザーの詳細情報をJSON形式で返す
	return c.JSON(fiber.Map{
		"id":         user.Id,
		"name":       user.Name,
		"bio":        user.Bio,
		"location":   user.Location,
		"website":    user.WibSite,
		"birth_date": user.Bio,
		"icon":       user.Icon,
		"banner":     user.Banner,
	})
}

// UpdateUser 関数は、ユーザー情報を更新するための関数です。
func UpdateUser(c *fiber.Ctx) error {
	// リクエストボディからデータを読み込むためのマップを定義
	var data map[string]string

	// リクエストボディの解析。エラーがあれば返す
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// ミドルウェアを通じて現在のユーザーIDを取得
	id, _ := middlewares.GetUserId(c)

	// 更新するユーザーのデータを取得するためのUserモデルのインスタンスを作成
	var user models.User

	// データベースからIDに基づいてユーザー情報を取得
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		// エラーがあれば、そのエラーを返す
		return result.Error
	}

	// ユーザーデータを更新
	user.Name = data["name"]
	user.Email = data["email"]
	user.Bio = data["bio"]
	user.Icon = data["icon"]
	user.Banner = data["banner"]
	user.Location = data["location"]
	user.WibSite = data["website"]
	user.BirthDate = data["birth_date"]

	// 更新されたデータをデータベースに保存
	database.DB.Save(&user)

	// 更新されたユーザー情報をJSON形式でレスポンスとして返す
	return c.JSON(user)
}

// UpdatePassword 関数は、ユーザーのパスワードを更新するための関数です。
func UpdatePassword(c *fiber.Ctx) error {
	// リクエストボディからデータを読み込むためのマップを定義
	var data map[string]string

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

	// ミドルウェアを通じて現在のユーザーIDを取得
	id, _ := middlewares.GetUserId(c)

	// 更新対象のユーザーモデルを作成
	user := models.User{
		Id: id,
	}

	// ユーザーモデルに新しいパスワードを設定
	user.SetPassword(data["password"])

	// データベースでユーザーの情報を更新
	database.DB.Model(&user).Updates(&user)

	// 更新されたユーザー情報をJSON形式でレスポンスとして返す
	return c.JSON(user)
}
