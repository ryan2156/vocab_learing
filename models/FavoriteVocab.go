package models

type FavoriteVocab struct {
	UserID       int    `json:"user_id"`
	VocabID      int    `json:"vocab_id"`
	JoinDate     string `json:"join_date"`
	ReadingCount int    `json:"reading_count"`
}
