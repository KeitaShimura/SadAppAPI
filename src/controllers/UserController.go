package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c *fiber.Ctx) error {
	// ページ番号とページサイズを取得
	page, pageSize := getPaginationParameters(c)

	var users []models.User
	var total int64
	database.DB.Model(&models.User{}).Count(&total)

	result := database.DB.
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&users)

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
		"email":      user.Email,
		"bio":        user.Bio,
		"location":   user.Location,
		"website":    user.WebSite,
		"birth_date": user.BirthDate,
		"icon":       user.Icon,
		"banner":     user.Banner,
	})
}

// UpdateUser 関数は、ユーザー情報を更新するための関数です。
func UpdateUser(c *fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	bio := c.FormValue("bio")
	location := c.FormValue("location")
	website := c.FormValue("website")
	birthDate := c.FormValue("birth_date")

	// 各フィールドのバリデーション
	if len(name) < 1 || len(name) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ユーザー名は1文字以上255文字以下である必要があります。",
		})
	}

	// ユーザーデータのバリデーション
	if len(name) < 1 || len(name) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ユーザー名は1文字以上255文字以下である必要があります。",
		})
	}

	if len(email) < 1 || len(email) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "メールアドレスは1文字以上255文字以下である必要があります。",
		})
	}
	if len(bio) > 1000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "自己紹介文は1000文字以下である必要があります。",
		})
	}

	if len(location) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "居住地のURは255文字以下である必要があります。",
		})
	}

	if len(website) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ウェブサイトのURLは255文字以下である必要があります。",
		})
	}

	if len(birthDate) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "生年月日は255文字以下である必要があります。",
		})
	}

	// ミドルウェアを通じて現在のユーザーIDを取得
	id, err := middlewares.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ユーザーIDの取得に失敗しました。",
		})
	}
	// 更新するユーザーのデータを取得するためのUserモデルのインスタンスを作成
	var user models.User

	// アイコン画像の取得と保存
	iconPath := ""
	iconFile, iconErr := c.FormFile("icon")
	if iconErr == nil {
		iconPath = filepath.Join("src/uploads", iconFile.Filename)
		if err := c.SaveFile(iconFile, iconPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "アイコン画像の保存に失敗しました。",
			})
		}
	}

	// バナー画像の取得と保存
	bannerPath := ""
	bannerFile, bannerErr := c.FormFile("banner")
	if bannerErr == nil {
		bannerPath = filepath.Join("src/uploads", bannerFile.Filename)
		if err := c.SaveFile(bannerFile, bannerPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "バナー画像の保存に失敗しました。",
			})
		}
	}

	// データベースからIDに基づいてユーザー情報を取得
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		// エラーがあれば、そのエラーを返す
		return result.Error
	}

	// リクエストデータから更新対象のユーザーデータを一時変数に格納
	updateData := map[string]interface{}{
		"Name":      name,
		"Email":     email,
		"Bio":       bio,
		"Location":  location,
		"WebSite":   website,
		"BirthDate": birthDate,
	}

	// Icon と Banner のパスを更新
	if iconPath != "" {
		updateData["Icon"] = "/" + iconPath
	}
	if bannerPath != "" {
		updateData["Banner"] = "/" + bannerPath
	}

	// 一時変数のデータをユーザーデータに反映
	database.DB.Model(&user).Updates(updateData)

	// 更新されたユーザー情報をJSON形式でレスポンスとして返す
	return c.JSON(user)
}

// UpdateEmail 関数は、ユーザー情報を更新するための関数です。
func UpdateEmail(c *fiber.Ctx) error {
	// リクエストボディからデータを読み込むための構造体を定義
	var input struct {
		Email string `json:"email"`
	}

	// リクエストボディの解析
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	if len(input.Email) < 1 || len(input.Email) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "メールアドレスは1文字以上255文字以下である必要があります。",
		})
	}

	// ミドルウェアを通じて現在のユーザーIDを取得
	id, _ := middlewares.GetUserId(c)

	// データベースからIDに基づいてユーザー情報を取得
	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		// エラーがあれば、そのエラーを返す
		return err
	}

	// 更新対象のユーザーデータを更新
	user.Email = input.Email
	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}

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

	// ミドルウェアを通じて現在のユーザーIDを取得
	id, _ := middlewares.GetUserId(c)

	// 更新対象のユーザーモデルを作成
	user := models.User{
		Id: id,
	}

	database.DB.Where("id = ?", id).First(&user)

	// 古いパスワードが一致しない場合、エラーを返す
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["current_password"])); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"error":   "current_password_mismatch",
			"message": "現在のパスワードが正しくありません。",
		})
	}

	// データベースからユーザー情報を取得
	if err := database.DB.Where(&user).First(&user).Error; err != nil {
		// ユーザーが見つからない場合や古いパスワードが一致しない場合はエラーレスポンスを返す
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "現在のパスワードが正しくありません。",
		})
	}

	// パスワードとパスワード確認が一致していない場合、400ステータスコードを返す。
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"error": "パスワードとパスワード確認が一致しません。",
		})
	}

	// パスワードのバリデーション：8文字以上
	if len(data["password"]) < 8 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "パスワードは8文字以上である必要があります。",
		})
	}

	// ユーザーモデルに新しいパスワードを設定
	user.SetPassword(data["password"])

	// データベースでユーザーの情報を更新
	database.DB.Model(&user).Updates(&user)

	// 更新されたユーザー情報をJSON形式でレスポンスとして返す
	return c.JSON(fiber.Map{
		"message": "パスワードが更新されました。",
	})
}
