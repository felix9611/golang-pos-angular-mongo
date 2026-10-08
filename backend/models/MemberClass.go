package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberClass struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ClassCode string             `bson:"classCode" json:"classCode"`
	ClassName string             `bson:"className" json:"className"`
	Remark    string             `bson:"remark" json:"remark"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}
