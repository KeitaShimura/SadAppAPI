package middlewares

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// IsAuthenticated はリクエストが認証済みかどうかを確認するミドルウェアです。
// JWT（Json Web Token）を使用して、ユーザーの認証状態を検証します。
func IsAuthenticated(c *fiber.Ctx) error {
	// Retrieve the JWT token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "認証されていません。",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Bearer token not found",
		})
	}

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	// Check for validity
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	// Proceed to the next handler
	return c.Next()
}

// GetUserId はリクエストからユーザーIDを取得するヘルパー関数です。
// JWT（Json Web Token）からユーザーの識別子を解析し、それを返します。
func GetUserId(c *fiber.Ctx) (uint, error) {
	// Authorization ヘッダーからトークンを取得
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("Authorization header is missing")
	}

	// Bearer トークンの形式を確認
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return 0, errors.New("Bearer token not found")
	}

	// JWT トークンを解析
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	// エラーがあれば0とエラーを返す
	if err != nil {
		return 0, err
	}

	// トークンの主題（Subject）からユーザーIDを解析
	payload := token.Claims.(*jwt.StandardClaims)
	id, _ := strconv.Atoi(payload.Subject)

	// ユーザーIDを返す
	return uint(id), nil
}
