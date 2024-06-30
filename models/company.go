package models

import "github.com/tyange/white-shadow-api/db"

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

func GetAllCompanyByUserId(userId *int64) ([]Company, error) {
	query := "SELECT * FROM companies WHERE user_id = ?"
	rows, err := db.DB.Query(query, userId)
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
