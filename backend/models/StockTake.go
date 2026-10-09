package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockTakes struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActionName    string             `bson:"actionName" json:"actionName"`
	ActionPlaceId string             `bson:"actionPlaceId" json:"actionPlaceId"`
	Remark        string             `bson:"remark" json:"remark"`
	Status        int                `bson:"status" json:"status"`
	CreatedTime   time.Time          `bson:"createdTime" json:"createdTime"`
	FinishTime    time.Time          `bson:"finishTime" json:"finishTime"`
	CreatedBy     string             `bson:"createdBy" json:"createdBy"`
	FinishBy      string             `bson:"finishBy" json:"finishBy"`
}

type StockTakeItems struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StockTakeId primitive.ObjectID `bson:"stockTakeId" json:"stockTakeId"`
	ProductId   primitive.ObjectID `bson:"productId" json:"productId"`
	ProductCode string             `bson:"productCode" json:"productCode"`
	PlaceId     string             `bson:"placeId" json:"placeId"`
	Qty         float64            `bson:"qty" json:"qty"`
	Status      string             `bson:"status" json:"status"`
	FinalStatus string             `bson:"finalStatus" json:"finalStatus"`
	CheckTime   time.Time          `bson:"checkTime,omitempty" json:"checkTime"`
	Remark      string             `bson:"remark" json:"remark"`
}
