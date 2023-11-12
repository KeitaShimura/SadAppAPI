package routes

import (
	"SadApp/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	// 'user' グループの下でルートを設定
	user := api.Group("user")
	// アカウント登録
	user.Post("register", controllers.Register)
}