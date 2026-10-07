package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Department struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DeptCode string             `bson:"deptCode" json:"deptCode"`
	DeptName string             `bson:"deptName" json:"deptName"`
	Remark   string             `bson:"remark" json:"remark"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}
