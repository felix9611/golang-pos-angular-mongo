package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Locations struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlaceCode      string             `bson:"placeCode" json:"placeCode"`
	PlaceName      string             `bson:"placeName" json:"placeName"`
	PlaceOtherName string             `bson:"placeOtherName" json:"placeOtherName"`
	Country        string             `bson:"country" json:"country"`
	Address        string             `bson:"address" json:"address"`
	ZipCode        string             `bson:"zipCode" json:"zipCode"`
	Email          string             `bson:"email" json:"email"`
	Phone          string             `bson:"phone" json:"phone"`
	Fax            string             `bson:"fax" json:"fax"`
	Remark         string             `bson:"remark" json:"remark"`
	Status         int                `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}
