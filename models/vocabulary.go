package models

type Vocabulary struct {
	VocabID     int    `json:"vocab_id"`
	Word        string `json:"word"`
	Defination  string `json:"defination"`
	Example_eng string `json:"example_eng"`
	Example_zh  string `json:"example_zh"`
	Part        int    `json:"part"`
	AddedDate   string `json:"added_date"`
	AddedBy     int    `json:"added_by"`
}

type Vocabulary_name struct {
	VocabID     int    `json:"vocab_id"`
	Word        string `json:"word"`
	Defination  string `json:"defination"`
	Example_eng string `json:"example_eng"`
	Example_zh  string `json:"example_zh"`
	Part        string `json:"part"`
	AddedDate   string `json:"added_date"`
	AddedBy     string `json:"added_by"`
}
