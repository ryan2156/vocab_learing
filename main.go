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
		AllowOrigins:     []string{"http://localhost:8080"}, // 允許的前端地址
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
		vocab.GET("/public", handlers.GetPublicVocabularies)            // 獲取公開單字庫 finished
		vocab.GET("/details/:vocab_id", handlers.GetVocabDetailHandler) // finished
		vocab.GET("/search", handlers.SearchPublicVocabHandler)

	}

	// 身份驗證相關路由
	authorized := r.Group("/")
	authorized.Use(handlers.AuthMiddleware())
	{
		authorized.GET("/profile", handlers.ProfileHandler)                                   // finished
		authorized.POST("/addVocab", handlers.AddVocabularyHandler)                           // 新增單字 finished / bug: 1 vocab add but two req get.
		authorized.POST("/addFavorite", handlers.AddFavoriteVocab)                            // 收藏最愛單字 finished
		authorized.GET("/vocabularies/added_by", handlers.GetAuthorVocabularies)              // finished 但是修改單字可以再優化
		authorized.PUT("/vocabularies/edit/:vocab_id", handlers.UpdateVocabulary)             // finished
		authorized.DELETE("/vocabularies/delete/:vocab_id", handlers.DeleteVocabularyHandler) // 刪除由自己新增的單字 // finished
		authorized.DELETE("/favorite/:vocab_id", handlers.RemoveFavoriteHandler)              // 移除收藏的單字
	}

	r.Run(":8888")
}
