package dto

import (
	"golang-fixedasset-mongo-backend/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockInOutProductLocationDto struct {
	ProductCode string `bson:"productCode" json:"productCode"`
	ProductId   string `bson:"productId" json:"productId"`
	PlaceCode   string `bson:"placeCode" json:"placeCode"`
	LocationId  string `bson:"locationId" json:"locationId"`
	Qty         int64  `bson:"qty" json:"qty"`
	TotalPrice  int64  `bson:"totalPrice" json:"totalPrice"`
	TotalCost   int64  `bson:"totalCost" json:"totalCost"`
}

type ListProductLocationtRequestDto struct {
	Page        int      `bson:"page" json:"page"`
	Limit       int      `bson:"limit" json:"limit"`
	LocationIds []string `bson:"locationIds" json:"locationIds"`
	ProductCode string   `bson:"productCode" json:"productCode"`
}

type ProductLocationLists struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	ProductId  string              `bson:"productId" json:"productId"`
	LocationId string              `bson:"locationId" json:"locationId"`
	Qty        int64               `bson:"qty" json:"qty"`
	TotalPrice int64               `bson:"totalPrice" json:"totalPrice"`
	TotalCost  int64               `bson:"totalCost" json:"totalCost"`
	Location   models.Locations    `bson:"location" json:"location"`
	Product    models.ProductLists `bson:"product" json:"product"`
}

type ListProductLocationts struct {
	Lists []ProductLocationLists `bson:"lists" json:"lists"`
	Total int64                  `bson:"total" json:"total"`
	Page  int                    `bson:"page" json:"page"`
	Limit int                    `bson:"limit" json:"limit"`
}
