package database

import "os"

var (
	DBUsername   = os.Getenv("DB_USERNAME")   // データベースのユーザー名
	DBPassword   = os.Getenv("DB_PASSWORD")   // データベースのパスワード
	DBHost       = os.Getenv("DB_HOST")       // データベースサーバーのホスト名またはIPアドレス
	DBPort       = os.Getenv("DB_PORT")       // データベースサーバーのポート
	DBName       = os.Getenv("DB_NAME")       // データベース名
	DBParameters = os.Getenv("DB_PARAMETERS") // データベース接続パラメーター
	Port         = os.Getenv("PORT") // データベース接続パラメーター
)
