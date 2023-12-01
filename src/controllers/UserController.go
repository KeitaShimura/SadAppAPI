package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "ユーザー情報の取得に失敗しました",
		})
	}
	return c.JSON(users)
}

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
		"website":    user.WebSite,
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

    // ユーザーデータのバリデーション
    if len(data["name"]) < 1 || len(data["name"]) > 255 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "ユーザー名は1文字以上255文字以下である必要があります。",
        })
    }

    if len(data["email"]) < 1 || len(data["email"]) > 255 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "メールアドレスは1文字以上255文字以下である必要があります。",
        })
    }

    if len(data["bio"]) > 1000 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "自己紹介文は1000文字以下である必要があります。",
        })
    }

    if len(data["website"]) > 255 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "ウェブサイトのURLは255文字以下である必要があります。",
        })
    }

    if len(data["birth_date"]) > 255 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "生年月日は255文字以下である必要があります。",
        })
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

    // リクエストデータから更新対象のユーザーデータを一時変数に格納
    updateData := map[string]interface{}{
        "Name":      data["name"],
        "Email":     data["email"],
        "Bio":       data["bio"],
        "Icon":      data["icon"],
        "Banner":    data["banner"],
        "Location":  data["location"],
        "WebSite":   data["website"],
        "BirthDate": data["birth_date"],
    }

    // 一時変数のデータをユーザーデータに反映
    database.DB.Model(&user).Updates(updateData)

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

	// パスワードのバリデーション：8文字以上
	if len(data["password"]) < 8 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "パスワードは8文字以上である必要があります。",
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
