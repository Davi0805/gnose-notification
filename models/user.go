package models

type User struct {
	ID        string   `json:"id"`
	CompanyIds []string `json:"companyIds"`
}