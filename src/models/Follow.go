package models

// User 構造体は、ユーザー情報を表します。
type Follow struct {
	Id          uint `json:"id"`
	FollowingId uint `json:"following_id"`
	FollowerId  uint `json:"follower_id"`
}
