package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvRecords struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AssetCode string             `bson:"assetCode,omitempty" json:"assetCode"`
	PlaceFrom string             `bson:"placeFrom,omitempty" json:"placeFrom"`
	PlaceTo   string             `bson:"placeTo,omitempty" json:"placeTo"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
}

type InventoryRecords struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProductId string             `bson:"ProductId,omitempty" json:"ProductId"`
	LocFrom   string             `bson:"locFrom,omitempty" json:"locFrom"`
	LocTo     string             `bson:"locTo,omitempty" json:"locTo"`
	Qty       int64              `bson:"qty" json:"qty"`
	Cost      int64              `bson:"cost" json:"cost"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
}
