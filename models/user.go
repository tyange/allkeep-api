package models

import (
	"errors"
	"time"

	"github.com/tyange/white-shadow-api/db"
	"github.com/tyange/white-shadow-api/utils"
)

type User struct {
	ID       int64
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Save() error {
	query := `INSERT INTO users(email, password, created_at) VALUES (?, ?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	currentTime := time.Now()
	result, err := stmt.Exec(u.Email, hashedPassword, currentTime)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	u.ID = userId

	return err
}

func (u *User) SaveWithoutPassword() error {
	query := `INSERT INTO users(email, created_at) VALUES (?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	currentTime := time.Now()
	result, err := stmt.Exec(u.Email, currentTime)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	u.ID = userId

	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil {
		return err
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}

func (u *User) CheckDuplicateUserId() bool {
	query := "SELECT id FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	err := row.Scan(&u.ID)

	return err == nil
}
