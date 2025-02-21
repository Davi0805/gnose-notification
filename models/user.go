package models

type User struct {
	ID        int   `json:"id"`
	CompanyIds []int `json:"companyIds"`
}