package dto

import (
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvRecordPureList struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	AssetCode     string             `bson:"assetCode,omitempty" json:"assetCode"`
	AssetName     string             `bson:"assetName,omitempty" json:"assetName"`
	PlaceFromCode string             `bson:"placeFromCode,omitempty" json:"placeFromCode"`
	PlaceFromName string             `bson:"placeFromName,omitempty" json:"placeFromName"`
	PlaceToCode   string             `bson:"placeToCode,omitempty" json:"placeToCode"`
	PlaceToName   string             `bson:"placeToName,omitempty" json:"placeToName"`
	CreatedAt     time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	PlaceFromData models.Locations   `bson:"placeFromData,omitempty" json:"placeFromData"`
	PlaceToData   models.Locations   `bson:"placeToData,omitempty" json:"placeToData"`
	AssetList     models.AssetLists  `bson:"assetList,omitempty" json:"assetList"`
}

type InvRecordList struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	AssetCode     string             `bson:"assetCode,omitempty" json:"assetCode"`
	PlaceFrom     string             `bson:"placeFrom,omitempty" json:"placeFrom"`
	PlaceTo       string             `bson:"placeTo,omitempty" json:"placeTo"`
	CreatedAt     time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	PlaceFromData models.Locations   `bson:"placeFromData,omitempty" json:"placeFromData"`
	PlaceToData   models.Locations   `bson:"placeToData,omitempty" json:"placeToData"`
	AssetList     models.AssetLists  `bson:"assetList,omitempty" json:"assetList"`
}

type InvRecordListResponse struct {
	Lists []InvRecordList `json:"lists"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

type InvRecordListDto struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}
