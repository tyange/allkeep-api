package models

import (
	"github.com/tyange/white-shadow-api/db"
	"github.com/tyange/white-shadow-api/utils"
)

type User struct {
	ID        int64
	Email     string `json:"email" binding:"required"`
	Passoword string `json:"password" binding:"required"`
}

func (u *User) Save() error {
	query := `INSERT INTO users(email, password) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Passoword)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	u.ID = userId

	return err
}
