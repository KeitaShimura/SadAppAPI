package controllers

import (
	"SadApp/src/database"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	// "golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupMockDB はモックされたデータベース接続をセットアップします
func setupMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	// GORMが接続時に実行するクエリを期待する
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("1"))

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when setting up the mock database connection", err)
	}

	return gormDB, mock
}

// setupRequest はテスト用のHTTPリクエストとレスポンスを準備します
func setupRequest(method, path string, body interface{}) (*fiber.App, *http.Request, *httptest.ResponseRecorder) {
	app := fiber.New()
	var req *http.Request

	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(bodyBytes))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	return app, req, res
}

// Register関数のテスト
func TestRegister(t *testing.T) {
	mockDB, _ := setupMockDB()
	database.DB = mockDB // モックデータベースをグローバルDBに設定
	// リクエストのセットアップ
	app, req, res := setupRequest("POST", "/api/user/register", map[string]string{
		"name":             "Test User",
		"emails":           "test@example.com",
		"password":         "password123",
		"password_confirm": "password123",
	})

	// ルートの登録
	app.Post("/api/user/register", Register)

	// リクエストの実行
	_, err := app.Test(req, -1)
	if err != nil {
		// エラー処理
		log.Printf("app.Test failed: %v", err)
		return
	}

	// ステータスコードの検証
	assert.Equal(t, http.StatusOK, res.Code) // 期待されるステータスコードを設定
}

// TestLogin 関数は、Login 関数のテストを行います。
func TestLogin(t *testing.T) {
	// モックデータベースのセットアップ
	mockDB, mock := setupMockDB()
	database.DB = mockDB

	// モックデータベースに存在するユーザー情報を設定
	mockedEmail := "keitashimura2023@gmail.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("11111111"), 12)
	mock.ExpectQuery("SELECT * FROM `users` WHERE email = ? ORDER BY `id` LIMIT 1").
		WithArgs(mockedEmail).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
			AddRow(1, "テスト1", mockedEmail, hashedPassword))

	// テスト用のリクエストデータ
	requestData := map[string]string{
		"email":    mockedEmail,
		"password": "11111111",
	}

	// Fiber アプリとテスト用リクエストをセットアップ
	app, req, res := setupRequest("POST", "/api/user/login", requestData)

	// ルートを登録
	app.Post("/api/user/login", Login)

	// リクエストをテスト
	_, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	// ステータスコードの検証
	assert.Equal(t, http.StatusOK, res.Code)
}
