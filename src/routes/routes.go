package routes

import (
	"SadApp/src/controllers"
	"SadApp/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// APIの基本パスを設定
	api := app.Group("/api")

	// ユーザー関連のルートを設定
	user := api.Group("/user")
	// アカウント登録
	user.Post("/register", controllers.Register)
	// ログイン
	user.Post("/login", controllers.Login)
	// ユーザー一覧を取得
	user.Get("/users", controllers.GetAllUsers)
	// 特定のユーザーの詳細を取得
	user.Get("/user/:id", controllers.GetUser)
	// 特定ユーザーのフォロワー一覧を取得
	user.Get("/followers/:id", controllers.GetFollowers)
	// 特定ユーザーがフォローしているユーザー一覧を取得
	user.Get("/followings/:id", controllers.GetFollowings)
	// 特定ユーザーの投稿一覧を取得
	user.Get("/user_posts/:id", controllers.UserPosts)
	// 特定ユーザーのイベント一覧を取得
	user.Get("/user_events/:id", controllers.UserEvents)

	// 認証が必要なユーザー関連のルート設定
	userAuthenticated := user.Use(middlewares.IsAuthenticated)
	// 認証済みユーザーの情報取得
	userAuthenticated.Get("/user", controllers.GetAuthUser)
	// ログアウト
	userAuthenticated.Post("/logout", controllers.Logout)
	// ユーザー情報の更新
	userAuthenticated.Put("/user", controllers.UpdateUser)
	// パスワードの更新
	userAuthenticated.Put("/user/password", controllers.UpdatePassword)

	// 投稿関連のルート設定
	posts := user.Group("/posts")
	// 投稿一覧取得
	posts.Get("", controllers.Posts)
	// 特定の投稿詳細取得
	posts.Get("/:id", controllers.GetPost)
	// 投稿のいいね数取得
	user.Get("/post/:id/likes", controllers.GetLikesForPost)

	// 認証が必要な投稿関連のルート設定
	userPostsAuthenticated := posts.Use(middlewares.IsAuthenticated)
	// 投稿の作成
	userPostsAuthenticated.Post("", controllers.CreatePost)
	// 投稿の更新
	userPostsAuthenticated.Put("/:id", controllers.UpdatePost)
	// 投稿の削除
	userPostsAuthenticated.Delete("/:id", controllers.DeletePost)
	// 投稿へのいいね
	userAuthenticated.Post("/post/:id/like", controllers.LikePost)
	// 投稿のいいねの解除
	userAuthenticated.Delete("/post/:id/unlike", controllers.UnlikePost)
	// 投稿がいいねされたかチェック
	userAuthenticated.Get("/post/:id/checklike", controllers.CheckIfPostLiked)

	// イベント関連のルート設定
	events := user.Group("/events")
	// イベント一覧取得
	events.Get("", controllers.Events)
	// 特定のイベント詳細取得
	events.Get("/:id", controllers.GetEvent)
	// イベントの参加者一覧
	events.Get("/:id/participants", controllers.GetEventParticipants)
	// イベントのいいね数取得
	user.Get("/event/:id/likes", controllers.GetLikesForEvent)

	// 認証が必要なイベント関連のルート設定
	userEventsAuthenticated := events.Use(middlewares.IsAuthenticated)
	// イベントの作成
	userEventsAuthenticated.Post("", controllers.CreateEvent)
	// イベントの更新
	userEventsAuthenticated.Put("/:id", controllers.UpdateEvent)
	// イベントの削除
	userEventsAuthenticated.Delete("/:id", controllers.DeleteEvent)
	// イベントへのいいね
	userAuthenticated.Post("/event/:id/like", controllers.LikeEvent)
	// イベントのいいねの解除
	userAuthenticated.Delete("/event/:id/unlike", controllers.UnlikeEvent)
	// イベントがいいねされたかチェック
	userAuthenticated.Get("/event/:id/checklike", controllers.CheckIfEventLiked)
	// イベントへの参加
	userEventsAuthenticated.Post("/:id/join", controllers.JoinEvent)
	// イベント参加の解除
	userEventsAuthenticated.Delete("/:id/leave", controllers.LeaveEvent)
}
