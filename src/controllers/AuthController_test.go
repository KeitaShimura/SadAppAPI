package controllers_test

import (
	"SadApp/src/controllers"
	"SadApp/src/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// モックデータベースをグローバル変数として定義
var mockDB *MockDB

func init() {
	mockDB = new(MockDB)
}

func TestRegister(t *testing.T) {
	// モックデータベースの設定
	mockDB.On("CreateUser", mock.Anything).Return(errors.New("データベースエラー"))

	// リクエストとレスポンスの設定
	app := fiber.New()
	app.Post("/api/user/register", controllers.Register)

	// テスト用データのエンコード
	body, _ := json.Marshal(map[string]string{
		"name":             "John Doe",
		"email":            "johndoe@example.com",
		"password":         "password123",
		"password_confirm": "password123",
	})

	// POSTリクエストの作成
	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// リクエストの送信
	resp, err := app.Test(req)

	// アサーション
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// モックの期待値チェック
	mockDB.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	// Fiberアプリケーションのセットアップ
	app := fiber.New()
	app.Post("/api/user/login", controllers.Login)

	// テスト用リクエストボディの準備
	reqBody := map[string]string{
		"email":    "1k1eitashimura202@gmail.com",
		"password": "11111111",
	}
	body, _ := json.Marshal(reqBody)

	// リクエストの作成
	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// リクエストの送信とレスポンスの取得
	resp, err := app.Test(req, -1)

	// アサーション
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 応答のボディを読み込むなど、追加のテストをここに記述
}

// 以下に MockDB と Database インターフェースを定義しています

type Database interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	// 他の必要なメソッドをここに追加
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}
