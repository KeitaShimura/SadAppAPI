package routes

import (
	"SadApp/src/controllers"
	"SadApp/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	// 'user' グループの下でルートを設定
	user := api.Group("user")
	// アカウント登録
	user.Post("register", controllers.Register)
	// ログイン
	user.Post("login", controllers.Login)
	// ユーザー一覧
	user.Get("users", controllers.GetAllUsers)
	// ユーザー詳細
	user.Get("user/:id", controllers.GetUser)
	// フォロワー一覧
	user.Get("followers/:id", controllers.GetFollowers)
	// フォローしているユーザー一覧
	user.Get("followings/:id", controllers.GetFollowings)

	// ユーザーごとの投稿一覧
	user.Get("user_posts/:id", controllers.UserPosts)

	// 'posts' グループの下でルートを設定
	posts := user.Group("posts")
	// 投稿一覧
	posts.Get("", controllers.Posts)

	// IsAuthenticatedミドルウェアを使用して、認証が必要なルートのグループを作成
	// このミドルウェアは、ユーザーが認証されているかどうかをチェックし、認証されていない場合は処理を進めない
	userAuthenticated := user.Use(middlewares.IsAuthenticated)
	// ユーザー認証
	userAuthenticated.Get("user", controllers.GetAuthUser)
	// ログアウト
	userAuthenticated.Post("logout", controllers.Logout)
	// ユーザー情報更新
	userAuthenticated.Put("user", controllers.UpdateUser)
	// パスワード更新
	userAuthenticated.Put("user/password", controllers.UpdatePassword)

	// フォロー
	userAuthenticated.Post("follow/:id", controllers.Follow)
	// フォロー解除
	userAuthenticated.Delete("unfollow/:id", controllers.UnFollow)

	// フォローチェック
	userAuthenticated.Get("check_if_following/:id", controllers.CheckIfFollowing)

	userPostsAuthenticated := posts.Use(middlewares.IsAuthenticated)

	// 投稿
	userPostsAuthenticated.Post("", controllers.CreatePost)
	// 投稿詳細取得
	posts.Get(":id", controllers.GetPost)
	// 投稿更新
	userPostsAuthenticated.Put(":id", controllers.UpdatePost)
	// 投稿削除
	userPostsAuthenticated.Delete(":id", controllers.DeletePost)

	// 'events' グループの下でルートを設定
	events := user.Group("events")
	userEventsAuthenticated := events.Use(middlewares.IsAuthenticated)

	// イベント一覧
	events.Get("", controllers.Events)
	// イベント
	userEventsAuthenticated.Post("", controllers.CreateEvent)
	// イベント詳細取得
	events.Get(":id", controllers.GetEvent)
	// イベント更新
	userEventsAuthenticated.Put(":id", controllers.UpdateEvent)
	// イベント削除
	userEventsAuthenticated.Delete(":id", controllers.DeleteEvent)
}
