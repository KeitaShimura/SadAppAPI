package controllers

import (
	"SadApp/src/database"
	// "SadApp/src/models"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	// "regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
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
        "email":            "test@example.com",
        "password":         "password123",
        "password_confirm": "password123",
    })

    // ルートの登録
    app.Post("/api/user/register", Register)

    // リクエストの実行
    app.Test(req, -1)

    // ステータスコードの検証
    assert.Equal(t, http.StatusOK, res.Code) // 期待されるステータスコードを設定
}

// TestLogin はLogin関数のテストです
// func TestLogin(t *testing.T) {
//     // モックデータベースのセットアップ
//     mockDB, mock := setupMockDB()
//     database.DB = mockDB

//     // モックが期待するクエリ
//     hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 14)
//     mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`email` = ? ORDER BY `users`.`id` LIMIT 1")).
//         WithArgs("test@example.com").
//         WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
//             AddRow(1, "Test User", "test@example.com", hashedPassword))


//     // テスト用のユーザーを作成
//     user := models.User{Name: "Test User", Email: "test@example.com"}
//     user.SetPassword("password123")

//     // データベースへのInsert操作をモック
// 	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`email` = ? ORDER BY `users`.`id` LIMIT 1")).
//     WithArgs("test@example.com").
//     WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
//         AddRow(1, "Test User", "test@example.com", "hashed_password"))

//     // Fiberアプリとテスト用リクエストをセットアップ
//     app, req, res := setupRequest("POST", "/api/user/login", map[string]string{
//         "email":    "test@example.com",
//         "password": "password123",
//     })

//     // ルートを登録
//     app.Post("/api/login", Login)

//     // リクエストをテスト
//     app.Test(req, -1)

//     // レスポンスのアサーション
//     assert.Equal(t, http.StatusOK, res.Code)

//     // モックの期待通りの動作を検証
//     if err := mock.ExpectationsWereMet(); err != nil {
//         t.Errorf("there were unfulfilled expectations: %s", err)
//     }
// }
