package dto

import "golang-fixedasset-mongo-backend/backend/models"

type ExcelFieldMatchPageDto struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
}

type ExcelFieldMatchList struct {
	Lists []models.ExcelFieldMatchs `json:"lists"`
	Total int64                     `json:"total"`
	Page  int                       `json:"page"`
	Limit int                       `json:"limit"`
}
