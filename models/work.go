package models

import (
	"time"

	"github.com/tyange/white-shadow-api/db"
)

type Work struct {
	ID          int64     `json:"id"`
	StartAt     time.Time `json:"start_at" binding:"required"`
	CompanyName string    `json:"company_name" binding:"required"`
	UserID      int64     `json:"user_id"`
}

func (w *Work) Save() error {
	query := `
	INSERT INTO works(start_at, company_name, user_id)
	VALUES (?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(w.StartAt, w.CompanyName, w.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	w.ID = id

	return err
}
