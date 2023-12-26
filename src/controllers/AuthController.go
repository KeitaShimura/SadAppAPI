package controllers

import (
	// 既存のインポート
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"net/mail"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt" // JWTを使用するためのパッケージをインポート
)

// ユーザー登録のための関数
func Register(c *fiber.Ctx) error {
	var data map[string]string

	// リクエストボディを解析する
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// 名前のバリデーション (3文字以上、20文字以下)
	nameLength := len(data["name"])
	if nameLength < 3 || nameLength > 20 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "名前は3文字以上20文字以下である必要があります。",
		})
	}

	// メールアドレスのバリデーション
	if !isValidEmail(data["email"]) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "無効なメールアドレスです。",
		})
	}

	// パスワードバリデーション (6文字以上)
	if len(data["password"]) < 6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "パスワードは6文字以上である必要があります。",
		})
	}

	// パスワードとパスワード確認の一致チェック
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "パスワードが一致しません。",
		})
	}

	// Userモデルのインスタンスを作成
	user := models.User{
		Name:  data["name"],
		Email: data["email"],
	}

	// パスワードのハッシュ化
	user.SetPassword(data["password"])

	// データベースにユーザーを保存
	database.DB.Create(&user)

	// JWTトークンの生成
	token, err := createJWTToken(user)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "トークンの生成に失敗しました。",
		})
	}

	// JWTトークンをクッキーに設定
	setTokenCookie(c, token)

	// レスポンスの準備
	response := prepareResponse(user, token)

	// レスポンスを返す
	return c.JSON(response)
}

// Login 関数は、ユーザーのログイン処理を行う関数です。
func Login(c *fiber.Ctx) error {
	var data map[string]string

	// リクエストボディの解析を試みる。
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// メールアドレスのバリデーション
	if !isValidEmail(data["email"]) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "無効なメールアドレスです。",
		})
	}

	// ユーザーモデルの新しいインスタンスを作成
	var user models.User

	// データベースからメールアドレスを使ってユーザー情報を検索する。
	database.DB.Where("email = ?", data["email"]).First(&user)

	// ユーザーが見つからなかった場合、エラーを返す。
	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "認証に失敗しました。",
		})
	}

	// パスワードが一致するかチェックする。一致しない場合はエラーを返す。
	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "認証に失敗しました。",
		})
	}

	// JWTトークンの生成
	token, err := createJWTToken(user)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "トークンの生成に失敗しました。",
		})
	}

	// クッキーにJWTトークンを設定する。
	setTokenCookie(c, token)

	// ユーザーデータにトークンを追加したJSONを作成
	response := prepareResponse(user, token)

	// 認証が成功したユーザーデータをJSON形式で返す。
	return c.JSON(response)
}

func GetAuthUser(c *fiber.Ctx) error {
	// ユーザーIDをミドルウェア関数から取得
	id, _ := middlewares.GetUserId(c)
	// ユーザーモデルの新しいインスタンスを作成
	var user models.User

	// データベースからユーザーIDに基づいてユーザー情報を取得
	database.DB.Where("id = ?", id).First(&user)

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

// メールアドレスの形式を確認するヘルパー関数
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// JWTトークンを生成する関数
func createJWTToken(user models.User) (string, error) {
	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))
	return token, err
}

// JWTトークンをクッキーに設定する関数
func setTokenCookie(c *fiber.Ctx, token string) {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "None",
		Domain:   "https://cocolo-talk.vercel.app", // ここにドメインを指定
	}
	c.Cookie(&cookie)
}


// レスポンスの準備をする関数
func prepareResponse(user models.User, token string) interface{} {
	return struct {
		Id        uint   `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Bio       string `json:"bio"`
		Icon      string `json:"icon"`
		Banner    string `json:"banner"`
		Location  string `json:"location"`
		WebSite   string `json:"website"`
		BirthDate string `json:"birth_date"`
		Token     string `json:"token"`
	}{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Bio:       user.Bio,
		Icon:      user.Icon,
		Banner:    user.Banner,
		Location:  user.Location,
		WebSite:   user.WebSite,
		BirthDate: user.BirthDate,
		Token:     token,
	}
}
