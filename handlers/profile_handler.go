package handlers

import (
	"fmt"
	"net/http"
	"vocab_learing/db"

	"github.com/gin-gonic/gin"
)

// ProfileHandler 返回用户资料和收藏单字
func ProfileHandler(c *gin.Context) {
	// 取得username
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 獲取使用者資料
	user, err := db.GetUserByUsername(username.(string))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	// 獲取使用者收藏的单字
	favorites, err := db.GetFavoriteVocabs(user.UserID)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorite vocabs"})
		return
	}

	fmt.Printf("Pass C\n")
	// 返回使用者資料和收藏的单字
	c.JSON(http.StatusOK, gin.H{
		"profile":   user,
		"favorites": favorites,
	})
}
