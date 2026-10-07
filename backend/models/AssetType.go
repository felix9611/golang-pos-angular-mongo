package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AssetTypes struct {
	ID       		 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TypeCode 		 string             `bson:"typeCode" json:"typeCode"`
	TypeName         string             `bson:"typeName" json:"typeName"`
	Remark           string             `bson:"remark" json:"remark"`
	DepreciationRate float64            `bson:"depreciationRate" json:"depreciationRate"`
	Status           int                `bson:"status" json:"status"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}