package router

import (
	"fmt"
	"os"
	"time"
	"w3st/factory"
	"w3st/infra"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() {

	// 環境変数 PORT を取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	db := infra.SetupDB()

	f := factory.NewFactory(db)


	// ルーターの初期化
	r := gin.Default()

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                      // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 許可するHTTPメソッド
		AllowHeaders:     []string{ "Access-Control-Allow-Credentials","Access-Control-Allow-Headers","Origin", "Content-Type", "Authorization"}, // 許可するヘッダー
		AllowCredentials: true,                                                // クレデンシャルを許可するかどうか
		MaxAge:           12 * time.Hour,                                      // キャッシュの最大時間
	}))

	// connectionTest
	r.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "connection success!!!!",
		})

	})

	// ユーザー
	users := r.Group("/users")
	userController := f.NewUserController()

	// ユーザー登録
	users.POST("/signup", userController.Signup)

	// ログイン
	users.POST("/login", userController.Login)

	// 指定されたポートでサーバーを開始
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}
