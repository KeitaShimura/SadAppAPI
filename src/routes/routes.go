package routes

import (
	"SadApp/src/controllers"
	"SadApp/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	// 'user' グループの下でルートを設定
	user := api.Group("user")
	// アカウント登録
	user.Post("register", controllers.Register)
	// ログイン
	user.Post("login", controllers.Login)
	// ユーザー詳細
	user.Get("user/:id", controllers.GetUser)

	userAuthenticated := user.Use(middlewares.IsAuthenticated)
	// ユーザー認証
	userAuthenticated.Get("user", controllers.GetAuthUser)
	// ログアウト
	userAuthenticated.Post("logout", controllers.Logout)
}
