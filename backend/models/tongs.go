package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tongs struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	Name      string             `bson:"name" json:"name"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}