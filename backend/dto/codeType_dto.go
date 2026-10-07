package dto

import "golang-fixedasset-mongo-backend/backend/models"

type CodeTypeListDto struct {
	Name  string `json:"name"`
	Limit int64  `json:"limit"`
	Page  int64  `json:"page"`
}

type CodeTypeList struct {
	Lists []models.CodeTypes `json:"lists"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}
