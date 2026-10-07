package dto

import "golang-fixedasset-mongo-backend/backend/models"

type TaxInformationListDTO struct {
	NameCode string `json:"nameCode"`
	Tax      string `json:"tax"`
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
}

type TaxInformationList struct {
	Lists []models.TaxInformations `json:"lists"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
}
