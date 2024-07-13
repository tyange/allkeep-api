package models

import (
	"fmt"
	"time"

	"github.com/tyange/white-shadow-api/db"
)

type Minutes int64

type Work struct {
	ID          int64     `json:"id"`
	CompanyID   int64     `json:"company_id" binding:"required"`
	CompanyName string    `json:"company_name" binding:"required"`
	WorkingTime Minutes   `json:"working_time" binding:"required"`
	StartAt     time.Time `json:"start_at"`
	DoneAt      time.Time `json:"done_at"`
	PauseAt     time.Time `json:"pause_at"`
	IsPause     bool      `json:"is_pause"`
	UserID      int64     `json:"user_id"`
}

type DuplicateCompanyIDError struct {
	CompanyID int64
}

func (e *DuplicateCompanyIDError) Error() string {
	return fmt.Sprintf("데이터가 이미 존재합니다: company_id %d", e.CompanyID)
}

func (w *Work) Save() error {
	checkQuery := "SELECT COUNT(*) FROM works WHERE company_id = ?"
	var count int
	err := db.DB.QueryRow(checkQuery, w.CompanyID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return &DuplicateCompanyIDError{CompanyID: w.CompanyID}
	}

	query := `
	INSERT INTO works(company_id, company_name, working_time, is_pause, user_id)
	VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(w.CompanyID, w.CompanyName, w.WorkingTime, w.IsPause, w.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	w.ID = id

	return err
}

func GetAllWorksByUserId(userId *int64) ([]Work, error) {
	query := "SELECT * FROM works WHERE user_id = ?"
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var works []Work
	var startAtString *string
	var doneAtSting *string
	var pauseAtString *string

	for rows.Next() {
		var work Work
		err := rows.Scan(&work.ID, &work.CompanyID, &work.CompanyName, &work.WorkingTime, &startAtString, &doneAtSting, &pauseAtString, &work.IsPause, &work.UserID)
		if err != nil {
			return nil, err
		}

		if startAtString != nil {
			work.StartAt, err = time.Parse("2006-01-02 15:04:05.000-07:00", *startAtString)
			if err != nil {
				return nil, err
			}
		}

		if doneAtSting != nil {
			work.DoneAt, err = time.Parse("2006-01-02 15:04:05.000-07:00", *doneAtSting)
			if err != nil {
				return nil, err
			}
		}

		if pauseAtString != nil {
			work.PauseAt, err = time.Parse("2006-01-02 15:04:05.000-07:00", *pauseAtString)
			if err != nil {
				return nil, err
			}
		}

		works = append(works, work)
	}

	return works, nil
}
