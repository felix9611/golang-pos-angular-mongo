package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExcelFieldMatchs struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FunctionCode string             `bson:"functionCode" json:"functionCode"`
	FunctionName string             `bson:"functionName" json:"functionName"`
	FunctionType string             `bson:"functionType" json:"functionType"`
	FieldLists   []FieldListsModel  `bson:"fieldLists" json:"fieldLists"`
	Status       int                `bson:"status" json:"status"`
	CreatedAt    time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

type FieldListsModel struct {
	DbFieldName    string `bson:"dbFieldName" json:"dbFieldName"`
	ExcelFieldName string `bson:"excelFieldName" json:"excelFieldName"`
	Sort           int    `bson:"sort" json:"sort"`
}
