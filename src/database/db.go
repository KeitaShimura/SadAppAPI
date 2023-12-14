package database

import (
	"SadApp/src/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB は、データベース接続を保持するためのグローバル変数です。
var DB *gorm.DB

// Connect はデータベースへの接続を確立する関数です。
func Connect() {
	var err error
	// 環境変数が設定されていない場合、デフォルトの情報を使用
	if DBUsername == "" {
		DBUsername = "root"
	}
	if DBPassword == "" {
		DBPassword = "@Keita8001"
	}
	if DBHost == "" {
		DBHost = "127.0.0.1"
	}
	if DBPort == "" {
		DBPort = "3306"
	}
	if DBName == "" {
		DBName = "sadapp"
	}
	if DBParameters == "" {
		DBParameters = "charset=utf8mb4&parseTime=True&loc=Local"
	}

	// dsn（データソース名）を組み立てます。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", DBUsername, DBPassword, DBHost, DBPort, DBName, DBParameters)

	// データベースに接続を開きます。
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode((logger.Silent)),
	})

	// データベース接続時のエラーをチェックします。
	if err != nil {
		// 接続に失敗した場合、プログラムをパニック状態にします。
		panic("データベースに接続できませんでした。") // "Could not connect to the database."
	}
}

// AutoMigrate はデータベースのスキーマを自動的にマイグレーション（更新）する関数です。
func AutoMigrate() {
	// Userモデルを使用して、データベースのスキーマを自動的にマイグレーションします。
	err := DB.AutoMigrate(models.User{}, models.Follow{}, models.Post{}, models.PostComment{}, models.PostLike{}, models.Event{}, models.EventComment{}, models.EventLike{}, models.EventParticipant{})
	if err != nil {
		// マイグレーションに失敗した場合、エラーをログに記録し、プログラムを終了します。
		log.Fatalf("データベースのマイグレーションに失敗しました: %v", err)
	}
}

type Database interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	// 他の必要なメソッドをここに追加
}

