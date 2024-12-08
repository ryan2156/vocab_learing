package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"vocab_learing/db"
	"vocab_learing/models"
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

// 獲取單字的詳細訊息
func GetVocabDetailHandler(c *gin.Context) {
	// 从 URL 中获取单词 ID
	vocabIDStr := c.Param("vocab_id")
	vocabID, err := strconv.Atoi(vocabIDStr)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vocab ID"})
		return
	}

	// 查询数据库获取单词详细信息
	vocabDetail, err := services.GetVocabDetailByID(vocabID)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vocab details"})
		return
	}

	// 如果没有找到单词，返回 404
	if vocabDetail == nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Vocab not found"})
		return
	}

	// 返回单词详细信息
	c.JSON(http.StatusOK, vocabDetail)
}

// 獲取由某使用者添加的單字庫
func GetAuthorVocabularies(c *gin.Context) {

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 調用 Service 層獲取單字庫
	vocabularies, err := services.GetAuthorVocabularies(username.(string))
	if err != nil {
		// 返回錯誤響應
		fmt.Printf("error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, vocabularies)
}

// 使用者修改由他自己創建的單字
func UpdateVocabulary(c *gin.Context) {
	var input struct {
		Word        string `json:"word"`
		Defination  string `json:"defination"`
		Example_eng string `json:"example_eng"`
		Example_zh  string `json:"example_zh"`
		Part        int    `json:"part"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	vocabIDStr := c.Param("vocab_id")
	vocabID, err := strconv.Atoi(vocabIDStr)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vocab ID"})
		return
	}

	var newVocab models.Vocabulary
	{
		newVocab.VocabID = vocabID
		newVocab.Word = input.Word
		newVocab.Defination = input.Defination
		newVocab.Example_eng = input.Example_eng
		newVocab.Example_zh = input.Example_zh
		newVocab.Part = input.Part
	}

	// 更新數據庫
	if err := db.UpdateVocabulary(vocabID, newVocab); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vocabulary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vocabulary updated successfully"})
}

// SearchPublicVocabHandler 搜尋公開單字
func SearchPublicVocabHandler(c *gin.Context) {
	// 獲取查詢參數
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供搜尋關鍵字"})
		return
	}

	// 執行 MSSQL 查詢
	query := `
		EXEC SearchPublicVocabularies @keyword = @p1
	`

	rows, err := db.DB.Query(query, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法搜尋單字庫", "details": err.Error()})
		return
	}
	defer rows.Close()

	// 解析結果

	var vocabularies []models.Vocabulary_name
	for rows.Next() {

		var vocab models.Vocabulary_name

		err = rows.Scan(
			&vocab.VocabID,
			&vocab.Word,
			&vocab.Defination,
			&vocab.Part,
			&vocab.AddedBy,
			&vocab.AddedDate,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析結果失敗", "details": err.Error()})
			return
		}

		vocabularies = append(vocabularies, vocab)
	}

	// 返回結果
	c.JSON(http.StatusOK, vocabularies)
}
