package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductLocations struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProductId  string             `bson:"productId" json:"productId"`
	LocationId string             `bson:"locationId" json:"locationId"`
	Qty        int64              `bson:"qty" json:"qty"`
	TotalPrice int64              `bson:"totalPrice" json:"totalPrice"`
	TotalCost  int64              `bson:"totalCost" json:"totalCost"`
}
