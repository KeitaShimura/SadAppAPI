package main

import (
	"SadApp/src/database"
	"SadApp/src/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	// データベースへの接続を確立します。
	database.Connect()
	// データベースのスキーマを自動マイグレーションします。
	database.AutoMigrate()

	// 新しいFiberアプリケーションのインスタンスを作成します。
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	// ルート設定をアプリケーションに追加します。
	routes.Setup(app)

	// サーバーをポート8002で起動します。
	// 何らかのエラーが発生した場合は、ログに記録してプログラムを終了します。
	log.Fatal(app.Listen(":8002"))
}
