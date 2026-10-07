package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WriteOffs struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssetId        string             `bson:"assetId" json:"assetId"`
	LastPlaceId    string             `bson:"lastPlaceId" json:"lastPlaceId"`
	Reason         string             `bson:"reason" json:"reason"`
	LastDay        time.Time          `bson:"lastDay,omitempty" json:"lastDay"`
	DisposalMethod string             `bson:"disposalMethod" json:"disposalMethod"`
	RemainingValue string             `bson:"remainingValue" json:"remainingValue"`
	Status         int                `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}
