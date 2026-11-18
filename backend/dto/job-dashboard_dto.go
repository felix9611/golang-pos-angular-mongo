package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type ListJobDashboardRequestDto struct {
	Name     string   `json:"name"`
	MemberId []string `json:"memberId"`
}

type UpdateJobDashboardStatusDto struct {
	ID        primitive.ObjectID `json:"_id"`
	Status    int                `json:"status"`
	SortOrder int                `json:"sortOrder"`
}
