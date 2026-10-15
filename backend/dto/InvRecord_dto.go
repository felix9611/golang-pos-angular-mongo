package dto

import (
	"golang-pos-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvRecordPureList struct {
	ProductId   string    `bson:"ProductId,omitempty" json:"ProductId"`
	LocFrom     string    `bson:"locFrom,omitempty" json:"locFrom"`
	LocTo       string    `bson:"locTo,omitempty" json:"locTo"`
	Qty         float64   `bson:"qty" json:"qty"`
	Cost        float64   `bson:"cost" json:"cost"`
	CreatedAt   time.Time `bson:"createdAt,omitempty" json:"createdAt"`
	ProductCode string    `bson:"productCode" json:"productCode"`
	ProductName string    `bson:"productName" json:"productName"`
	Unit        string    `bson:"unit" json:"unit"`
	ItemCode    string    `bson:"itemCode" json:"itemCode"`
	BrandCode   string    `bson:"brandCode" json:"brandCode"`
	BrandName   string    `bson:"brandName" json:"brandName"`
	CostPrice   float64   `bson:"costPrice" json:"costPrice"`
	RetailPrice float64   `bson:"retailPrice" json:"retailPrice"`
	LocFromCode string    `bson:"locFromCode" json:"locFromCode"`
	LocFromName string    `bson:"locFromName" json:"locFromName"`
	LocToCode   string    `bson:"locToCode" json:"locToCode"`
	LocToName   string    `bson:"locToName" json:"locToName"`
}

type InvRecordList struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	ProductId   string              `bson:"ProductId,omitempty" json:"ProductId"`
	LocFrom     string              `bson:"locFrom,omitempty" json:"locFrom"`
	LocTo       string              `bson:"locTo,omitempty" json:"locTo"`
	Qty         float64             `bson:"qty" json:"qty"`
	Cost        float64             `bson:"cost" json:"cost"`
	CreatedAt   time.Time           `bson:"createdAt,omitempty" json:"createdAt"`
	Product     models.ProductLists `bson:"product" json:"product"`
	LocFromData models.Locations    `bson:"locFromData" json:"locFromData"`
	LocToData   models.Locations    `bson:"locToData" json:"locToData"`
}

type InvRecordListResponse struct {
	Lists []InvRecordList `json:"lists"`
	Total int64           `json:"total"`
	Page  int64           `json:"page"`
	Limit int64           `json:"limit"`
}

type InvRecordListDto struct {
	Date  []string `json:"date"`
	Limit int64    `json:"limit"`
	Page  int64    `json:"page"`
}

type InsertInvRecordDto struct {
	ProductId string  `bson:"ProductId,omitempty" json:"ProductId"`
	LocFrom   string  `bson:"locFrom,omitempty" json:"locFrom"`
	LocTo     string  `bson:"locTo,omitempty" json:"locTo"`
	Qty       float64 `bson:"qty" json:"qty"`
	Cost      float64 `bson:"cost" json:"cost"`
}
