package handlers

import (
	"fmt"
	"net/http"
	"vocab_learing/db"
	"vocab_learing/services"

	"github.com/gin-gonic/gin"
)

// 添加新單字
func AddVocabularyHandler(c *gin.Context) {
	var input struct {
		Word        string `json:"word"`
		Defination  string `json:"defination"`
		Example_eng string `json:"example_eng"`
		Example_zh  string `json:"example_zh"`
		Part        int    `json:"part"`
	}

	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	adder_name := c.GetString("username")
	adder, err := db.GetUserByUsername(adder_name)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	adder_id := adder.UserID

	err = services.AddVocabulary(input.Word, input.Defination, input.Example_eng, input.Example_zh, input.Part, adder_id)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add vocabulary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vocabulary added successfully"})
}

// 獲取公開單字庫
func GetPublicVocabularies(c *gin.Context) {
	// 調用 Service 層獲取單字庫
	vocabularies, err := services.GetPublicVocabularies()
	if err != nil {
		// 返回錯誤響應
		fmt.Printf("error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, vocabularies)
}

// 加入最喜歡的單字
func AddFavoriteVocab(c *gin.Context) {
	var input struct {
		VocabID   int    `json:"vocab_id"`
		AddedDate string `json:"added_date"`
	}

	// 查看POST資料是否正常
	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 邏輯是： 透過username 查找 user_id，之後將(user_id, vocab_id, added_date)丟到函數run
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	err := services.AddFavoriteVocab(username, input.VocabID, input.AddedDate)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add favorite vocabulary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"succes": "succes"})
}
