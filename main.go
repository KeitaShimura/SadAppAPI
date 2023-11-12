package main

import (
	"SadApp/src/database"
	"SadApp/src/routes"
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
	// ルート設定をアプリケーションに追加します。
	routes.Setup(app)

	// サーバーをポート8002で起動します。
	// 何らかのエラーが発生した場合は、ログに記録してプログラムを終了します。
	log.Fatal(app.Listen(":8002"))
}
