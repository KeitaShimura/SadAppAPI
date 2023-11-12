package middlewares

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// IsAuthenticated はリクエストが認証済みかどうかを確認するミドルウェアです。
// JWT（Json Web Token）を使用して、ユーザーの認証状態を検証します。
func IsAuthenticated(c *fiber.Ctx) error {
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

	// JWTが有効な場合、次のミドルウェアまたはリクエストハンドラーに処理を渡す
	// c.Next() は、このミドルウェアの後に定義されている次の処理へと進むことを意味する
	return c.Next()
}

// GetUserId はリクエストからユーザーIDを取得するヘルパー関数です。
// JWT（Json Web Token）からユーザーの識別子を解析し、それを返します。
func GetUserId(c *fiber.Ctx) (uint, error) {
	// ユーザーのブラウザから"jwt"という名前のCookieを取得
	cookie := c.Cookies("jwt")

	// 取得したCookieを使用してJWTを解析
	// ここでは"secret"をキーとして使用している
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	// エラーがあれば0とエラーを返す
	if err != nil {
		return 0, err
	}

	// トークンの主題（Subject）からユーザーIDを解析
	payload := token.Claims.(*jwt.StandardClaims)

	// トークンの主題（Subject）からユーザーIDを解析
	id, _ := strconv.Atoi(payload.Subject)

	// ユーザーIDを返す
	return uint(id), nil
}
