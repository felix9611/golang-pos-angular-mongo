package dto

import (
	"golang-pos-mongo-backend/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductLocationPureList struct {
	ProductCode string  `bson:"productCode" json:"productCode"`
	ProductName string  `bson:"productName" json:"productName"`
	ItemCode    string  `bson:"itemCode" json:"itemCode"`
	BrandCode   string  `bson:"brandCode" json:"brandCode"`
	BrandName   string  `bson:"brandName" json:"brandName"`
	CostPrice   float64 `bson:"costPrice" json:"costPrice"`
	RetailPrice float64 `bson:"retailPrice" json:"retailPrice"`
	Unit        string  `bson:"unit" json:"unit"`
	PlaceCode   string  `bson:"placeCode" json:"placeCode"`
	PlaceName   string  `bson:"placeName" json:"placeName"`
	Qty         float64 `bson:"qty" json:"qty"`
	TotalPrice  float64 `bson:"totalPrice" json:"totalPrice"`
	TotalCost   float64 `bson:"totalCost" json:"totalCost"`
}

type StockInOutProductLocationDto struct {
	ProductCode string  `bson:"productCode" json:"productCode"`
	ProductId   string  `bson:"productId" json:"productId"`
	PlaceCode   string  `bson:"placeCode" json:"placeCode"`
	LocationId  string  `bson:"locationId" json:"locationId"`
	Qty         float64 `bson:"qty" json:"qty"`
	TotalPrice  float64 `bson:"totalPrice" json:"totalPrice"`
	TotalCost   float64 `bson:"totalCost" json:"totalCost"`
}

type StockMoveProductLocationDto struct {
	ProductCode    string  `bson:"productCode" json:"productCode"`
	ProductId      string  `bson:"productId" json:"productId"`
	FromPlaceCode  string  `bson:"fromPlaceCode" json:"fromPlaceCode"`
	FromLocationId string  `bson:"fromLocationId" json:"fromLocationId"`
	ToPlaceCode    string  `bson:"toPlaceCode" json:"toPlaceCode"`
	ToLocationId   string  `bson:"toLocationId" json:"toLocationId"`
	Qty            float64 `bson:"qty" json:"qty"`
	TotalPrice     float64 `bson:"totalPrice" json:"totalPrice"`
	TotalCost      float64 `bson:"totalCost" json:"totalCost"`
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
	Qty        float64             `bson:"qty" json:"qty"`
	TotalPrice float64             `bson:"totalPrice" json:"totalPrice"`
	TotalCost  float64             `bson:"totalCost" json:"totalCost"`
	Location   models.Locations    `bson:"location" json:"location"`
	Product    models.ProductLists `bson:"product" json:"product"`
}

type ListProductLocationts struct {
	Lists []ProductLocationLists `bson:"lists" json:"lists"`
	Total int64                  `bson:"total" json:"total"`
	Page  int                    `bson:"page" json:"page"`
	Limit int                    `bson:"limit" json:"limit"`
}

type CheckProductAndLocationDto struct {
	ProductId  string `bson:"productId" json:"productId"`
	LocationId string `bson:"locationId" json:"locationId"`
}

type CheckProductAndLocationResponse struct {
	Status bool                    `bson:"status" json:"status"`
	Data   models.ProductLocations `bson:"data" json:"data"`
}
