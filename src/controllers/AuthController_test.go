package controllers_test

import (
	"SadApp/src/controllers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// Fiberアプリケーションのセットアップ
	app := fiber.New()

	fmt.Println(app);
	// ルートの設定
	app.Post("/api/user/register", controllers.Register)

	// テスト用ユーザーデータの作成
	user := map[string]string{
		"name":             "志村",
		"email":            "1k1eitashimura202@gmail.com",
		"password":         "11111111",
		"password_confirm": "11111111",
	}

	// データのJSONエンコード
	body, _ := json.Marshal(user)

	// POSTリクエストの作成
	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// リクエストの送信とレスポンスの取得
	resp, err := app.Test(req)

	// エラーがないことの確認
	assert.Nil(t, err)

	// ステータスコードの確認
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestLogin(t *testing.T) {
	// Fiberアプリケーションのセットアップ
	app := fiber.New()

	// ルートの設定
	app.Post("/api/user/login", controllers.Login)

	// テスト用ログインデータの作成
	user := map[string]string{
		"email":    "1k1eitashimura202@gmail.com",
		"password": "11111111",
	}

	// データのJSONエンコード
	body, _ := json.Marshal(user)

	// POSTリクエストの作成
	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// リクエストの送信とレスポンスの取得
	resp, err := app.Test(req)

	// エラーがないことの確認
	assert.Nil(t, err)

	// ステータスコードの確認
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 応答のボディを読み込む
	// 例：レスポンスに含まれる特定のキーの存在や値をチェックするなど
	// response := make(map[string]interface{})
	// json.NewDecoder(resp.Body).Decode(&response)
	// assert.Equal(t, "期待される値", response["特定のキー"])
}
