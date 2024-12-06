package models

type User struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	PswdHash string `json:"pswd_hash"`
	JoinDate string `json:"join_date"`
	Email    string `json:"email"`
}
