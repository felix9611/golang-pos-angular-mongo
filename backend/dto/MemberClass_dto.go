package dto

import "golang-pos-mongo-backend/backend/models"

type MemberClassListDto struct {
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
	Name  string `json:"name"`
}

type MemberClassList struct {
	Lists []models.MemberClass `json:"lists"`
	Total int64                `json:"total"`
	Page  int64                `json:"page"`
	Limit int64                `json:"limit"`
}
