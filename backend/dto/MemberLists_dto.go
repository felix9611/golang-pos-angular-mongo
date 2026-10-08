package dto

import (
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberListsDto struct {
	Name  string `json:"name"`
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
}

type MemberListsByNameDto struct {
	Name string `json:"name"`
}

type MemberListDetail struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	MemberCode  string             `json:"memberCode" bson:"memberCode"`
	Name        string             `json:"name" bson:"name"`
	Address     string             `json:"address" bson:"address"`
	Email       string             `json:"email" bson:"email"`
	Phone       string             `json:"phone" bson:"phone"`
	Fax         string             `json:"fax" bson:"fax"`
	ClassId     string             `json:"classId" bson:"classId"`
	Remark      string             `json:"remark" bson:"remark"`
	Status      int                `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
	MemberClass models.MemberClass `json:"memberClass"`
	ClassName   string             `json:"className" bson:"className"`
	ClassCode   string             `json:"classCode" bson:"classCode"`
}

type MemberList struct {
	Lists []MemberListDetail `json:"lists"`
	Total int64              `json:"total"`
	Page  int64              `json:"page"`
	Limit int64              `json:"limit"`
}
