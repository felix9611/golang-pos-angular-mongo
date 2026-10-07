package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type InvRecords struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AssetCode    string             `bson:"assetCode,omitempty" json:"assetCode"`
	PlaceFrom   string             `bson:"placeFrom,omitempty" json:"placeFrom"`
	PlaceTo     string             `bson:"placeTo,omitempty" json:"placeTo"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
}