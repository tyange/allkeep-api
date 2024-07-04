package models

import (
	"fmt"

	"github.com/tyange/white-shadow-api/db"
)

type Company struct {
	ID          int64
	CompanyName string `json:"company_name"`
	UserID      int64  `json:"user_id"`
}

func (c *Company) Save() error {
	query := `INSERT INTO companies(company_name, user_id) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(c.CompanyName, c.UserID)

	if err != nil {
		return err
	}

	companyId, err := result.LastInsertId()

	c.ID = companyId

	return err
}

func GetAllCompanyByUserId(userId *int64, pageSize *int64, pageNum *int64) ([]Company, error) {
	offset := (*pageNum - 1) * *pageSize
	fmt.Println(offset)
	query := "SELECT * FROM companies WHERE user_id = ? LIMIT ? OFFSET ?"
	rows, err := db.DB.Query(query, userId, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company

	for rows.Next() {
		var company Company
		err := rows.Scan(&company.ID, &company.CompanyName, &company.UserID)
		if err != nil {
			return nil, err
		}

		companies = append(companies, company)
	}

	return companies, nil
}

func GetCompanyCountByUserId(userId *int64) (int64, error) {
	query := "SELECT COUNT(*) FROM companies WHERE user_id = ?"
	var count int64
	err := db.DB.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetCompanyById(userId *int64) (*Company, error) {
	query := "SELECT * FROM companies WHERE id = ?"
	row := db.DB.QueryRow(query, userId)

	var company Company
	err := row.Scan(&company.ID, &company.CompanyName, &company.UserID)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (company Company) Update() error {
	query := `
	UPDATE companies
	SET company_name = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(company.CompanyName, company.ID)
	return err
}
