package models

import (
	"time"

	"github.com/tyange/white-shadow-api/db"
)

type Company struct {
	ID          int64      `json:"id"`
	CompanyName string     `json:"company_name"`
	UserID      int64      `json:"user_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type AllCompanyResponseData struct {
	Companies      []Company `json:"companies"`
	TotalItemCount int64     `json:"total_item_count"`
}

func (c *Company) Save() error {
	query := `INSERT INTO companies(company_name, user_id, created_at) VALUES (?, ?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	currentTime := time.Now()
	result, err := stmt.Exec(c.CompanyName, c.UserID, currentTime)

	if err != nil {
		return err
	}

	companyId, err := result.LastInsertId()

	c.ID = companyId

	return err
}

func GetCompaniesByUserIdWithPage(userId *int64, pageSize *int64, pageNum *int64) (*AllCompanyResponseData, error) {
	offset := (*pageNum - 1) * *pageSize
	companiesQuery := "SELECT * FROM companies WHERE user_id = ? LIMIT ? OFFSET ?"
	rows, err := db.DB.Query(companiesQuery, userId, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	countQuery := "SELECT COUNT(*) FROM companies WHERE user_id = ?"
	var count int64
	err = db.DB.QueryRow(countQuery, userId).Scan(&count)
	if err != nil {
		return nil, err
	}

	var data AllCompanyResponseData

	data.TotalItemCount = count

	for rows.Next() {
		var company Company
		err := rows.Scan(&company.ID, &company.CompanyName, &company.UserID, &company.CreatedAt, &company.UpdatedAt)
		if err != nil {
			return nil, err
		}

		data.Companies = append(data.Companies, company)
	}

	return &data, nil
}

func GetCompaniesByUserId(userId *int64) ([]Company, error) {
	query := "SELECT * FROM companies WHERE user_id = ?"
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company

	for rows.Next() {
		var company Company
		err := rows.Scan(&company.ID, &company.CompanyName, &company.UserID, &company.CreatedAt, &company.UpdatedAt)
		if err != nil {
			return nil, err
		}

		companies = append(companies, company)
	}

	return companies, nil
}

func GetCompanyById(companyId *int64) (*Company, error) {
	query := `SELECT * FROM companies WHERE id = ?`
	row := db.DB.QueryRow(query, companyId)

	var company Company
	err := row.Scan(&company.ID, &company.CompanyName, &company.UserID, &company.CreatedAt, &company.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (company Company) Update() error {
	query := `
	UPDATE companies
	SET company_name = ?, updated_at = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	currentTime := time.Now()
	_, err = stmt.Exec(company.CompanyName, currentTime, company.ID)
	return err
}

func (company Company) Delete() error {
	query := "DELETE FROM companies WHERE id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(company.ID)

	return err
}
