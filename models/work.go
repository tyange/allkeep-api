package models

import (
	"fmt"
	"time"

	"github.com/tyange/white-shadow-api/db"
)

type Minutes int64

type Work struct {
	ID          int64      `json:"id"`
	CompanyID   int64      `json:"company_id" binding:"required"`
	CompanyName string     `json:"company_name" binding:"required"`
	WorkingTime Minutes    `json:"working_time" binding:"required"`
	StartAt     *time.Time `json:"start_at"`
	DoneAt      *time.Time `json:"done_at"`
	PauseAt     *time.Time `json:"pause_at"`
	IsPause     bool       `json:"is_pause"`
	IsDone      bool       `json:"is_done"`
	UserID      int64      `json:"user_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
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
	INSERT INTO works(company_id, company_name, working_time, is_pause, is_done, user_id, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	currentTime := time.Now()
	result, err := stmt.Exec(w.CompanyID, w.CompanyName, w.WorkingTime, w.IsPause, w.IsDone, w.UserID, currentTime)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	w.ID = id

	return err
}

func GetAllWorksByUserId(userId *int64) ([]Work, error) {
	query := "SELECT * FROM works WHERE user_id = ? ORDER BY created_at DESC"
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var works []Work

	for rows.Next() {
		var work Work
		err := rows.Scan(&work.ID, &work.CompanyID, &work.CompanyName, &work.WorkingTime, &work.StartAt, &work.DoneAt, &work.PauseAt, &work.IsPause, &work.IsDone, &work.UserID, &work.CreatedAt, &work.UpdatedAt)
		if err != nil {
			return nil, err
		}

		works = append(works, work)
	}

	return works, nil
}

func GetWorkById(workId *int64) (*Work, error) {
	query := `SELECT * FROM works WHERE id = ?`
	row := db.DB.QueryRow(query, workId)

	var work Work
	err := row.Scan(&work.ID, &work.CompanyID, &work.CompanyName, &work.WorkingTime, &work.StartAt, &work.DoneAt, &work.PauseAt, &work.IsPause, &work.IsDone, &work.UserID, &work.CreatedAt, &work.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &work, nil
}

func UpdateWorkForStart(workId *int64, startAt *time.Time, doneAt *time.Time) error {
	query := `
	UPDATE works
	SET start_at = ?,
		done_at = ?,
		updated_at = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	currentTime := time.Now()
	_, err = stmt.Exec(startAt, doneAt, currentTime, &workId)

	return err
}

func UpdateWorkForPause(workId *int64, pauseAt *time.Time) error {
	query := `
	UPDATE works
	SET pause_at = ?,
		is_pause = ?,
		updated_at = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	currentTime := time.Now()
	_, err = stmt.Exec(pauseAt, true, currentTime, &workId)

	return err
}

func UpdateWorkForRestart(workId *int64, doneAt *time.Time) error {
	query := `
	UPDATE works
	SET done_at = ?,
		is_pause = ?,
		updated_at = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	currentTime := time.Now()
	_, err = stmt.Exec(doneAt, false, currentTime, &workId)

	return err
}

func UpdateWorkForDone(workId *int64) error {
	query := `
	UPDATE works
	SET is_done = ?,
		updated_at = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	currentTime := time.Now()
	_, err = stmt.Exec(true, currentTime, &workId)

	return err
}

func (work Work) Update() error {
	query := `
	UPDATE works
	SET company_id = ?,
		company_name = ?,
		working_time = ?,
		start_at = ?,
		done_at = ?,
		pause_at = ?,
		is_pause = ?,
		user_id = ?,
		updated_at = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	currentTime := time.Now()
	_, err = stmt.Exec(work.CompanyID, work.CompanyName, work.WorkingTime, work.StartAt, work.DoneAt, work.PauseAt, work.IsPause, work.UserID, currentTime, work.ID)

	return err
}
