package controllers

import (
	// 既存のインポート
	"SadApp/src/database"
	"SadApp/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt" // JWTを使用するためのパッケージをインポート
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

	// Userモデルのインスタンスを作成する。
	user := models.User{
		Name:  data["name"],
		Email: data["email"],
	}

	// ユーザーのパスワードをハッシュ化して保存
	user.SetPassword(data["password"])

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

// Logout はユーザーのログアウト処理を行う関数です。
func Logout(c *fiber.Ctx) error {
	// 新しいクッキーを作成し、名前を"jwt"に設定
	// 値を空にし、有効期限を過去の時間に設定することでクッキーを削除
	cookie := fiber.Cookie{
		Name:     "jwt",                      // クッキーの名前
		Value:    "",                         // クッキーの値を空に設定
		Expires:  time.Now().Add(-time.Hour), // 有効期限を過去に設定してクッキーを無効化
		HTTPOnly: true,                       // JavaScriptからのアクセスを防ぐための設定
	}

	// 設定したクッキーをレスポンスに追加
	c.Cookie(&cookie)

	// ログアウト成功のメッセージをJSON形式で返す
	return c.JSON(fiber.Map{
		"message": "ログアウトしました。",
	})
}
