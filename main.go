package main

import (
	"SadApp/src/database"
	"SadApp/src/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
		AllowOrigins:     "http://localhost:3000, https://cocolo-talk.vercel.app/", // フロントエンドのオリジンを具体的に指定
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

	// ルート設定をアプリケーションに追加します。
	routes.Setup(app)

	// サーバーをポート8002で起動します。
	// 何らかのエラーが発生した場合は、ログに記録してプログラムを終了します。
	port := database.Port
	if port == "" {
		port = "8080" // デフォルトポート
	}
	log.Fatal(app.Listen(":" + port))

}