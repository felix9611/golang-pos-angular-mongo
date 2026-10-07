package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SysMenus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	MainId    string             `bson:"mainId" json:"mainId"`
	Name      string             `bson:"name" json:"name"`
	Icon      string             `bson:"icon" json:"icon"`
	Path      string             `bson:"path" json:"path"`
	Sort      int                `bson:"sort" json:"sort"`
	Type      string            `bson:"type" json:"type"`
	ExcelFunctionCode string             `bson:"excelFunctionCode" json:"excelFunctionCode"`
	ExcelFunctionName string             `bson:"excelFunctionName" json:"excelFunctionName"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}