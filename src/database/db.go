package database

import (
	"SadApp/src/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	// 接続文字列を設定ファイルから組み立てる
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", DBUsername, DBPassword, DBHost, DBPort, DBName, DBParameters)

	// Open the connection to the database
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Check for errors in opening the database connection
	if err != nil {
		panic("データベースに接続できませんでした。") // "Could not connect to the database."
	}
}

func AutoMigrate() {
	DB.AutoMigrate(models.User{})
}
