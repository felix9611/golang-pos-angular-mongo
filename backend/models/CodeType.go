package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CodeTypes struct {
	ID   	  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ValueCode string             `bson:"valueCode" json:"valueCode"`
	ValueName string             `bson:"valueName" json:"valueName"`
	Type 	  string             `bson:"type" json:"type"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}