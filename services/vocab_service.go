package services

import (
	"fmt"
	"vocab_learing/db"
	"vocab_learing/models"
)

// SearchVocabByWord 根據單字查詢資料庫中的Vocabularis檔案
func SearchVocabByWord(word string) ([]models.Vocabulary, error) {
	// Select ，有用Like 來查詢相關的單字
	query := `
		SELECT v.vocab_id, v.word, v.defination, p.name AS part_name, u.username AS added_by
		FROM Vocabularies v
		JOIN Parts p ON v.part = p.part_id
		JOIN Users u ON v.added_by = u.user_id
		WHERE v.word LIKE @p1
	`

	// 執行 Select
	rows, err := db.DB.Query(query, "%"+word+"%") // 模糊匹配
	if err != nil {
		return nil, fmt.Errorf("failed to search vocab by word: %v", err)
	}
	defer rows.Close()

	// 遍歷撈出來的資料並整理
	var vocabularies []models.Vocabulary
	for rows.Next() {
		var vocab models.Vocabulary
		if err := rows.Scan(&vocab.VocabID, &vocab.Word, &vocab.Defination, &vocab.Part, &vocab.AddedBy); err != nil {
			return nil, fmt.Errorf("failed to scan vocab row: %v", err)
		}
		vocabularies = append(vocabularies, vocab)
	}

	// 看Scan的錯誤
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return vocabularies, nil
}

func AddVocabulary(word, defination, example_eng, example_zh string, part, addedBy int) error {
	query := `INSERT INTO Vocabularies (word, defination, example_eng, example_zh, part, added_by) VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`
	_, err := db.DB.Exec(query, word, defination, example_eng, example_zh, part, addedBy)
	return err
}

func AddFavoriteVocab(username any, vocabID int, added_Date string) error {

	user, err := db.GetUserByUsername(username.(string))
	if err != nil {
		return err
	}

	query := `INSERT INTO Favorite_Vocabs (user_id, vocab_id, join_date) VALUES (@p1, @p2, @p3)`
	_, err = db.DB.Exec(query, user.UserID, vocabID, added_Date)
	return err
}

func GetPublicVocabularies() ([]models.Vocabulary_name, error) {
	// 查询公开单字库，同时获取用户名
	query := `
		SELECT v.vocab_id, v.word, v.defination, p.name AS part_name, u.username AS added_by_name
		FROM Vocabularies v
		JOIN Users u ON v.added_by = u.user_id
		JOIN Parts p ON v.part = p.part_id
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 解析数据库结果
	var vocabularies []models.Vocabulary_name
	for rows.Next() {
		var vocab models.Vocabulary_name
		if err := rows.Scan(&vocab.VocabID, &vocab.Word, &vocab.Defination, &vocab.Part, &vocab.AddedBy); err != nil {
			return nil, err
		}
		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, nil
}

func GetVocabDetails() {

}
