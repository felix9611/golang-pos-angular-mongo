package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Locations struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlaceCode string             `bson:"placeCode" json:"placeCode"`
	PlaceName string             `bson:"placeName" json:"placeName"`
	Remark    string             `bson:"remark" json:"remark"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}