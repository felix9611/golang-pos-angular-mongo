package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryRecords struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProductId string             `bson:"ProductId,omitempty" json:"ProductId"`
	LocFrom   string             `bson:"locFrom,omitempty" json:"locFrom"`
	LocTo     string             `bson:"locTo,omitempty" json:"locTo"`
	Qty       float64            `bson:"qty" json:"qty"`
	Cost      float64            `bson:"cost" json:"cost"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
}
