package database

import (
	"SadApp/src/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// DB は、データベース接続を保持するためのグローバル変数です。
var DB *gorm.DB

// Connect はデータベースへの接続を確立する関数です。
func Connect() {
	var err error
	// dsn（データソース名）を組み立てます。これには接続情報が含まれます。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", DBUsername, DBPassword, DBHost, DBPort, DBName, DBParameters)

	// データベースに接続を開きます。
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// データベース接続時のエラーをチェックします。
	if err != nil {
		// 接続に失敗した場合、プログラムをパニック状態にします。
		panic("データベースに接続できませんでした。") // "Could not connect to the database."
	}
}

// AutoMigrate はデータベースのスキーマを自動的にマイグレーション（更新）する関数です。
func AutoMigrate() {
	// Userモデルを使用して、データベースのスキーマを自動的にマイグレーションします。
	err := DB.AutoMigrate(models.User{}, models.Follow{})
	if err != nil {
		// マイグレーションに失敗した場合、エラーをログに記録し、プログラムを終了します。
		log.Fatalf("データベースのマイグレーションに失敗しました: %v", err)
	}
}
