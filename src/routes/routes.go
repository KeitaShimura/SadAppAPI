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
	// 投稿の更新（認証が必要）
	userPostsAuthenticated.Put("/:id", controllers.UpdatePost)
	// 投稿の削除（認証が必要）
	userPostsAuthenticated.Delete("/:id", controllers.DeletePost)
	// 投稿へのいいね
	userAuthenticated.Post("/post/:id/like", controllers.LikePost)
	// 投稿のいいねの解除
	userAuthenticated.Delete("/post/:id/unlike", controllers.UnlikePost)
	// 投稿がいいねされたかチェック
	userAuthenticated.Get("/post/:id/checklike", controllers.CheckIfPostLiked)

	// コメント関連のルート設定（投稿）
	postComments := api.Group("/comments/post")
	// 特定の投稿に対するコメント一覧を取得
	postComments.Get("/:post_id", controllers.PostComments)
	// コメントの作成（認証が必要）
	userPostCommentsAuthenticated := postComments.Use(middlewares.IsAuthenticated)
	userPostCommentsAuthenticated.Post("/", controllers.CreatePostComment)
	// コメントの更新（認証が必要）
	userPostCommentsAuthenticated.Put("/:id", controllers.UpdatePostComment)
	// コメントの削除（認証が必要）
	userPostCommentsAuthenticated.Delete("/:id", controllers.DeletePostComment)

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
	// イベントの作成（認証が必要）
	userEventsAuthenticated.Post("", controllers.CreateEvent)
	// イベントの更新（認証が必要）
	userEventsAuthenticated.Put("/:id", controllers.UpdateEvent)
	// イベントの削除（認証が必要）
	userEventsAuthenticated.Delete("/:id", controllers.DeleteEvent)
	// イベントへのいいね（認証が必要）
	userAuthenticated.Post("/event/:id/like", controllers.LikeEvent)
	// イベントのいいねの解除（認証が必要）
	userAuthenticated.Delete("/event/:id/unlike", controllers.UnlikeEvent)
	// イベントがいいねされたかチェック（認証が必要）
	userAuthenticated.Get("/event/:id/checklike", controllers.CheckIfEventLiked)
	// イベントへの参加（認証が必要）
	userEventsAuthenticated.Post("/:id/join", controllers.JoinEvent)
	// イベント参加の解除（認証が必要）
	userEventsAuthenticated.Delete("/:id/leave", controllers.LeaveEvent)

	// コメント関連のルート設定（イベント）
	eventComments := api.Group("/comments/event")
	// 特定のイベントに対するコメント一覧を取得
	eventComments.Get("/:event_id", controllers.EventComments)
	// コメントの作成（認証が必要）
	userEventCommentsAuthenticated := eventComments.Use(middlewares.IsAuthenticated)
	userEventCommentsAuthenticated.Post("/", controllers.CreateEventComment)
	// コメントの更新（認証が必要）
	userEventCommentsAuthenticated.Put("/:id", controllers.UpdateEventComment)
	// コメントの削除（認証が必要）
	userEventCommentsAuthenticated.Delete("/:id", controllers.DeleteEventComment)
}
