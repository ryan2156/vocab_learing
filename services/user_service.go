package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"vocab_learing/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO Users (username, pswd_hash, email) VALUES (@p1, @p2, @p3)`
	_, err = db.DB.Exec(query, username, string(hashedPassword), email)
	return err
}

func AuthenticateUser(username, password string) (bool, error) {
	var hashedPassword string
	query := `SELECT pswd_hash FROM Users WHERE username = @p1`
	err := db.DB.QueryRow(query, username).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		fmt.Printf("err: %s\n", err)
		return false, nil // 用戶不存在
	} else if err != nil {
		fmt.Printf("err: %s\n", err)
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, nil
}

func GetUserProfile(c *gin.Context) {
	username := c.GetString("username")

	userinfo, err := db.GetUserByUsername(username)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userinfo": userinfo,
	})
}
