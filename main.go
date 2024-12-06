package main

import (
	"time"
	"vocab_learing/db"
	"vocab_learing/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.InitDB(); err != nil {
		panic(err)
	}

	r := gin.Default()

	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://10.221.3.165:8080"}, // 允許的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // 是否允許跨域攜帶 Cookie
		MaxAge:           12 * time.Hour, // OPTIONS 請求的緩存時間
	}))

	r.POST("/register", handlers.RegisterHandler) // finished
	r.POST("/login", handlers.LoginHandler)       // finished

	// 單字庫相關路由
	vocab := r.Group("/vocabularies")
	{
		vocab.GET("/public", handlers.GetPublicVocabularies) // 獲取公開單字庫 finished

	}

	// 身份驗證相關路由
	authorized := r.Group("/")
	authorized.Use(handlers.AuthMiddleware())
	{
		authorized.GET("/profile", handlers.ProfileHandler)         // finished
		authorized.POST("/addVocab", handlers.AddVocabularyHandler) // 新增單字 finished
		authorized.POST("/addFavorite", handlers.AddFavoriteVocab)  // 收藏最愛單字
	}

	r.Run(":8888")
}
