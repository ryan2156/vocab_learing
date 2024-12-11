package db

import (
	"database/sql"
	"errors"
	"fmt"
	"vocab_learing/models"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

func InitDB() error {
	var err error

	// 配置 SQL Server 連接信息
	server := "localhost"
	port := 1433 // 默認端口
	user := "sa"
	password := "Mysuperpswd456"
	database := "vocab_learing"

	// 創建連接字符串
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		user, password, server, port, database)

	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		return err
	}

	// 驗證連接是否成功
	err = DB.Ping()
	if err != nil {
		return err
	}

	// 初始化表格
	queries := []string{
		// 1. Users 表
		`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'Users' AND xtype = 'U')
		CREATE TABLE Users (
			user_id INT IDENTITY(1,1) PRIMARY KEY,
			name NVARCHAR(30) NOT NULL UNIQUE,
			pswd_hash NVARCHAR(255) NOT NULL,
			join_date DATE DEFAULT GETDATE(),
			email NVARCHAR(100) NOT NULL UNIQUE
		);
		`,
		// 2. Parts 表
		`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'Parts' AND xtype = 'U')
		CREATE TABLE Parts (
			part_id INT IDENTITY(1,1) PRIMARY KEY,
			name NVARCHAR(20) NOT NULL
		);
		`,
		// 3. Vocabularies 表
		`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'Vocabularies' AND xtype = 'U')
		CREATE TABLE Vocabularies (
			vocab_id INT IDENTITY(1,1) PRIMARY KEY,
			word NVARCHAR(70) NOT NULL,
			defination NTEXT NOT NULL,
			example_eng NTEXT,
			example_zh NTEXT,
			part INT NOT NULL,
			join_date DATE DEFAULT GETDATE(),
			added_by INT NOT NULL,
			FOREIGN KEY (part) REFERENCES Parts(part_id) ON DELETE CASCADE,
			FOREIGN KEY (added_by) REFERENCES Users(user_id) ON DELETE CASCADE
		);
		`,
		// 4. Favorite_Vocabs 表
		`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'Favorite_Vocabs' AND xtype = 'U')
		CREATE TABLE Favorite_Vocabs (
			user_id INT NOT NULL,
			vocab_id INT NOT NULL,
			join_date DATE DEFAULT GETDATE(),
			reading_count INT DEFAULT 0,
			PRIMARY KEY (user_id, vocab_id),
			FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
			FOREIGN KEY (vocab_id) REFERENCES Vocabularies(vocab_id) ON DELETE CASCADE
		);
		`,
	}

	// 執行所有表格創建語句
	for _, query := range queries {
		_, err = DB.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

// 從使用者名稱獲取資料
func GetUserByUsername(username string) (*models.User, error) {
	// 定義要返回的 User 結構體
	var user models.User

	// SQL 查詢語句
	query := `SELECT user_id, username, pswd_hash, email, join_date FROM Users WHERE username = @p1`

	// 執行查詢
	row := DB.QueryRow(query, username)

	// 將查詢結果掃描到 User 結構體中
	err := row.Scan(&user.UserID, &user.Name, &user.PswdHash, &user.Email, &user.JoinDate)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果沒有找到用戶
			return nil, fmt.Errorf("user not found")
		}
		// 如果是其他數據庫錯誤
		return nil, err
	}

	// 返回用戶信息
	return &user, nil
}

func GetFavoriteVocabs(userID int) ([]models.Vocabulary_name, error) {
	query := `
		SELECT v.vocab_id, v.word, v.defination, v.example_eng, v.example_zh, p.name AS part_name, v.added_date
		FROM Favorite_Vocabs fv
		JOIN Vocabularies v ON fv.vocab_id = v.vocab_id
		JOIN Parts p ON v.part = p.part_id
		WHERE fv.user_id = @p1
	`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []models.Vocabulary_name
	for rows.Next() {
		var vocab models.Vocabulary_name
		if err := rows.Scan(&vocab.VocabID, &vocab.Word, &vocab.Defination, &vocab.Example_eng, &vocab.Example_zh, &vocab.Part, &vocab.AddedDate); err != nil {
			return nil, err
		}
		favorites = append(favorites, vocab)
	}

	return favorites, nil
}

func UpdateVocabulary(vocabID int, vocab models.Vocabulary) error {
	query := "EXEC UpdateVocabulary @VocabID, @Word, @Defination, @ExampleEng, @ExampleZh, @Part"
	_, err := DB.Exec(query,
		sql.Named("VocabID", vocabID),
		sql.Named("Word", vocab.Word),
		sql.Named("Defination", vocab.Defination),
		sql.Named("ExampleEng", vocab.Example_eng),
		sql.Named("ExampleZh", vocab.Example_zh),
		sql.Named("Part", vocab.Part),
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteVocabulary(userID, vocabID int) error {
	// 構建查詢語句，檢查單字是否屬於當前用戶
	query := `
		DELETE FROM Vocabularies
		WHERE vocab_id = @p1 AND added_by = @p2
	`

	// 執行刪除操作
	result, err := DB.Exec(query, vocabID, userID)
	if err != nil {
		return err
	}

	// 檢查刪除是否成功
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		ErrForbidden := errors.New("無權移除單字")
		return ErrForbidden // 無權刪除，或者單字不存在
	}

	return nil
}

var ErrFavoriteNotFound = errors.New("favorite not found")

func RemoveFavorite(userID, vocabID int) error {
	query := `
        DELETE FROM Favorite_Vocabs
        WHERE user_id = @p1 AND vocab_id = @p2
    `
	result, err := DB.Exec(query, userID, vocabID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrFavoriteNotFound // 如果沒有刪除任何記錄
	}

	return nil
}
