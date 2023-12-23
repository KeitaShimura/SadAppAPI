package database

import (
	"SadApp/src/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 環境変数をロードする関数
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// DB は、データベース接続を保持するためのグローバル変数です。
var DB *gorm.DB

// Connect はデータベースへの接続を確立する関数です。
func Connect() {
	loadEnv()

	DBUsername := os.Getenv("DB_USERNAME")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	DBName := os.Getenv("DB_NAME")
	DBParameters := os.Getenv("DB_PARAMETERS")

	var err error

	// PostgreSQL用のDSNを組み立てます。
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s", DBHost, DBUsername, DBPassword, DBName, DBPort, DBParameters)

	// データベースに接続を開きます。
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	// データベース接続時のエラーをチェックします。
	if err != nil {
		panic("データベースに接続できませんでした。")
	}
}

// AutoMigrate はデータベースのスキーマを自動的にマイグレーション（更新）する関数です。
func AutoMigrate() {
	// Userモデルを使用して、データベースのスキーマを自動的にマイグレーションします。
	err := DB.AutoMigrate(&models.User{}, &models.Follow{}, &models.Post{}, &models.PostComment{}, &models.PostLike{}, &models.Event{}, &models.EventComment{}, &models.EventLike{}, &models.EventParticipant{})
	if err != nil {
		log.Fatalf("データベースのマイグレーションに失敗しました: %v", err)
	}
}
