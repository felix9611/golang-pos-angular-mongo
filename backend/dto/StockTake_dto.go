package dto

import (
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListStockTakeDto struct {
	PlaceIds []string `json:"placeIds"`
	Name     string   `json:"name"`
	Page     int64    `json:"page"`
	Limit    int64    `json:"limit"`
}

type GetStockTakeResponse struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActionName    string             `bson:"actionName" json:"actionName"`
	ActionPlaceId string             `bson:"actionPlaceId" json:"actionPlaceId"`
	Remark        string             `bson:"remark" json:"remark"`
	Status        int                `bson:"status" json:"status"`
	CreatedTime   time.Time          `bson:"createdTime" json:"createdTime"`
	FinishTime    *time.Time         `bson:"finishTime" json:"finishTime"`
	CreatedBy     string             `bson:"createdBy" json:"createdBy"`
	FinishBy      string             `bson:"finishBy" json:"finishBy"`
	Location      models.Locations   `bson:"location" json:"location"`
}

type StockTakeResponseList struct {
	Lists []GetStockTakeResponse `json:"lists"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	Limit int                    `json:"limit"`
}

type GetStockTakeItemResponse struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StockTakeId primitive.ObjectID `bson:"stockTakeId" json:"stockTakeId"`
	AssetId     primitive.ObjectID `bson:"assetId" json:"assetId"`
	AssetCode   string             `bson:"assetCode" json:"assetCode"`
	PlaceId     string             `bson:"placeId" json:"placeId"`
	Status      string             `bson:"status" json:"status"`
	CheckTime   time.Time          `bson:"checkTime,omitempty" json:"checkTime"`
	Remark      string             `bson:"remark" json:"remark"`
	Location    models.Locations   `bson:"location" json:"location"`
	Assetlist   AssetListsDetails  `bson:"assetlist" json:"assetlist"`
}

type GetStockTakeFormWithItemsResponse struct {
	ID             primitive.ObjectID         `bson:"_id,omitempty" json:"id"`
	ActionName     string                     `bson:"actionName" json:"actionName"`
	ActionPlaceId  string                     `bson:"actionPlaceId" json:"actionPlaceId"`
	Remark         string                     `bson:"remark" json:"remark"`
	Status         int                        `bson:"status" json:"status"`
	CreatedTime    time.Time                  `bson:"createdTime" json:"createdTime"`
	FinishTime     time.Time                  `bson:"finishTime" json:"finishTime"`
	CreatedBy      string                     `bson:"createdBy" json:"createdBy"`
	FinishBy       string                     `bson:"finishBy" json:"finishBy"`
	StockTakeItems []GetStockTakeItemResponse `bson:"stockTakeItems" json:"stockTakeItems"`
}

/*
return gin.H{
		"id":             stockTakeForm.ID,
		"actionName":     stockTakeForm.ActionName,
		"actionPlaceId":  stockTakeForm.ActionPlaceId,
		"remark":         stockTakeForm.Remark,
		"createdTime":    stockTakeForm.CreatedTime,
		"status":         stockTakeForm.Status,
		"createdBy":      stockTakeForm.CreatedBy,
		"finishTime":     stockTakeForm.FinishTime,
		"finishBy":       stockTakeForm.FinishBy,
		"stockTakeItems": stockTakeItems,
	}, nil
*/
