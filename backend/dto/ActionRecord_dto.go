package dto

import "golang-fixedasset-mongo-backend/backend/models"

type ListActionRecordReqDto struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ActionReocrdList struct {
	Lists []models.ActionRecords `json:"lists"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	Limit int                    `json:"limit"`
}
