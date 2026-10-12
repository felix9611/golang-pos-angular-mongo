package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberLists struct {
	Id                primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	MemberCode        string              `json:"memberCode" bson:"memberCode"`
	Name              string              `json:"name" bson:"name"`
	Address           string              `json:"address" bson:"address"`
	Email             string              `json:"email" bson:"email"`
	Phone             string              `json:"phone" bson:"phone"`
	Fax               string              `json:"fax" bson:"fax"`
	ClassId           string              `json:"classId" bson:"classId"`
	Remark            string              `json:"remark" bson:"remark"`
	Status            int                 `bson:"status" json:"status"`
	CreatedAt         time.Time           `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt         time.Time           `bson:"updatedAt,omitempty" json:"updatedAt"`
	MemberSpecialDays []MemberSpecialDays `json:"memberSpecialDays" bson:"memberSpecialDays"`
}

type MemberSpecialDays struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	MemberId  string             `json:"memberId" bson:"memberId"`
	Name      string             `json:"name" bson:"name"`
	Date      string             `json:"date" bson:"date"`
	Remark    string             `json:"remark" bson:"remark"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}
