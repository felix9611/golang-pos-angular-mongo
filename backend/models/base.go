package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Base struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}