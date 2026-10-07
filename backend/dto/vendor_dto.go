package dto

import "golang-fixedasset-mongo-backend/backend/models"

type VendorPageDto struct {
	Page    int64  `json:"page"`
	Limit   int64  `json:"limit"`
	Name    string `json:"name"`
	Place   string `json:"place"`
	Contact string `json:"contact"`
}

type VendorList struct {
	Lists []models.Vendors `json:"lists"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}
