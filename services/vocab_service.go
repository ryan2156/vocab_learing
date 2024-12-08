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
	// 查询 PublicVocabularies 视图
	query := "SELECT vocab_id, word, defination, part_name, added_by_name, added_date FROM VW_PublicVocabs"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vocabularies []models.Vocabulary_name
	for rows.Next() {
		var vocab models.Vocabulary_name
		err := rows.Scan(
			&vocab.VocabID,
			&vocab.Word,
			&vocab.Defination,
			&vocab.Part,
			&vocab.AddedBy,
			&vocab.AddedDate,
		)
		if err != nil {
			return nil, err
		}
		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, nil
}

func GetAuthorVocabularies(author string) ([]models.Vocabulary_name, error) {
	query := `
		SELECT 
			vocab_id,
			word,
			defination,
			part_name,
			added_by_name, 
			added_date 
		FROM VW_PublicVocabs
		WHERE added_by_name = @p1
	`
	rows, err := db.DB.Query(query, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vocabularies []models.Vocabulary_name
	for rows.Next() {
		var vocab models.Vocabulary_name
		err := rows.Scan(
			&vocab.VocabID,
			&vocab.Word,
			&vocab.Defination,
			&vocab.Part,
			&vocab.AddedBy,
			&vocab.AddedDate,
		)
		if err != nil {
			return nil, err
		}
		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, nil
}

func GetVocabDetailByID(vocabID int) (*models.Vocabulary_name, error) {
	query := `
		SELECT vocab_id, word, definition, example_eng, example_zh, part_name, added_by_name, added_date 
		FROM VocabDetails 
		WHERE vocab_id = @p1
	`
	row := db.DB.QueryRow(query, vocabID)

	var detail models.Vocabulary_name
	err := row.Scan(
		&detail.VocabID,
		&detail.Word,
		&detail.Defination,
		&detail.Example_eng,
		&detail.Example_zh,
		&detail.Part,
		&detail.AddedBy,
		&detail.AddedDate,
	)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}
