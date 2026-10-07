package dto

import "golang-fixedasset-mongo-backend/backend/models"

type ProductTypeListDto struct {
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
	Name  string `json:"name"`
}

type ProductTypeList struct {
	Lists []models.ProductTypes `json:"lists"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}
