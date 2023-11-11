package main

import (
	"SadApp/src/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// データベースへの接続を確立します。
	database.Connect()
	// データベースのスキーマを自動マイグレーションします。
	database.AutoMigrate()

	// 新しいFiberアプリケーションのインスタンスを作成します。
	app := fiber.New()

	// ルートパス ('/') に対するルートを定義します。
	app.Get("/", func(c *fiber.Ctx) error {
		// リクエストに対して "Hello, World!" というレスポンスを返します。
		return c.SendString("Hello, World!")
	})

	// サーバーをポート8002で起動します。
	// 何らかのエラーが発生した場合は、ログに記録してプログラムを終了します。
	log.Fatal(app.Listen(":8002"))
}
