package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductTypes struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TypeCode      string             `bson:"typeCode" json:"typeCode"`
	TypeName      string             `bson:"typeName" json:"typeName"`
	TypeOtherName string             `bson:"typeOtherName" json:"typeOtherName"` // typeOtherName
	Remark        string             `bson:"remark" json:"remark"`
	Status        int                `bson:"status" json:"status"`
	CreatedAt     time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}
