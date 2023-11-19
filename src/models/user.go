package models

import "golang.org/x/crypto/bcrypt"

// User 構造体は、ユーザー情報を表します。
type User struct {
	Id        uint   `json:"id"`                               // ユーザーのID
	Name      string `json:"name"`                             // ユーザーの名前
	Email     string `json:"email" gorm:"unique"`              // ユーザーのメールアドレス
	Password  []byte `json:"-"`                                // ハッシュ化されたパスワード
	Bio       string `json:"bio" gorm:"nullable"`              // ユーザーの自己紹介文
	Icon      string `json:"icon" gorm:"type:text;nullable"`   // ユーザーのアイコン画像
	Banner    string `json:"banner" gorm:"type:text;nullable"` // ユーザーの背景画像
	Location  string `json:"location" gorm:"nullable"`         // ユーザーの居住地
	WibSite   string `json:"website" gorm:"nullable"`          // ユーザーのウェブサイト
	BirthDate string `json:"birth_date" gorm:"nullable"`       // ユーザーの生年月日
}

// SetPassword は、与えられたパスワードをハッシュ化してUser構造体に設定します。
func (user *User) SetPassword(password string) {
	// bcryptを使用してパスワードをハッシュ化
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = hashedPassword // ハッシュ化されたパスワードをUser構造体に設定
}

// ComparePassword は、与えられたパスワードがユーザーのパスワードと一致するかをチェックします。
func (user *User) ComparePassword(password string) error {
	// bcryptを使用してハッシュ化されたパスワードと比較
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
