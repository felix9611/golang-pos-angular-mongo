package dto

import "golang-pos-mongo-backend/backend/models"

type ExcelFieldMatchPageDto struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type ExcelFieldMatchList struct {
	Lists []models.ExcelFieldMatchs `json:"lists"`
	Total int                       `json:"total"`
	Page  int                       `json:"page"`
	Limit int                       `json:"limit"`
}
