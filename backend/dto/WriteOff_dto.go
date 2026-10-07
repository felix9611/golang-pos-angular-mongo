package dto

import (
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateWriteOffRecrod struct {
	AssetId        string  `json:"assetId"`
	LastPlaceId    string  `json:"lastPlaceId"`
	Reason         string  `json:"reason"`
	LastDay        string  `json:"lastDay"`
	DisposalMethod string  `json:"disposalMethod"`
	RemainingValue float64 `json:"remainingValue"`
}

type ListWriteOffReqDto struct {
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	PlaceIds  []string    `json:"placeIds"`
	DeptIds   []string    `json:"deptIds"`
	TypeIds   []string    `json:"typeIds"`
	DateRange []time.Time `json:"dateRange"`
}

type WriteOffPureList struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssetCode      string             `bson:"assetCode" json:"assetCode"`
	AssetName      string             `bson:"assetName" json:"assetName"`
	PurchaseDate   string             `bson:"purchaseDate" json:"purchaseDate"`
	LastPlaceCode  string             `bson:"lastPlaceCode" json:"lastPlaceCode"`
	LastPlaceName  string             `bson:"lastPlaceName" json:"lastPlaceName"`
	Reason         string             `bson:"reason" json:"reason"`
	LastDay        time.Time          `bson:"lastDay,omitempty" json:"lastDay"`
	DisposalMethod string             `bson:"disposalMethod" json:"disposalMethod"`
	RemainingValue string             `bson:"remainingValue" json:"remainingValue"`
	Status         int                `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
	AssetList      models.AssetLists  `bson:"assetlist,omitempty" json:"assetlist"`
	Location       models.Locations   `bson:"location,omitempty" json:"location"`
}

type WriteOffList struct {
	Lists []models.WriteOffs `json:"lists"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}
