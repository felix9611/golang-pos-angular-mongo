package dto

import "golang-fixedasset-mongo-backend/backend/models"

type DepartmentPageDto struct {
	Name  string `json:"name"`
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
}

type DepartmentList struct {
	Lists []models.Department `json:"lists"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}
