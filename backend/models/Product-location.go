package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductLocations struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProductId  string             `bson:"productId" json:"productId"`
	LocationId string             `bson:"locationId" json:"locationId"`
	Qty        float64            `bson:"qty" json:"qty"`
	TotalPrice float64            `bson:"totalPrice" json:"totalPrice"`
	TotalCost  float64            `bson:"totalCost" json:"totalCost"`
}
