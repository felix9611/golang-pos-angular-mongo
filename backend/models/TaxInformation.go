package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaxInformations struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NationCode string             `bson:"nationCode" json:"nationCode"`
	NationName string             `bson:"nationName" json:"nationName"`
	CountryCode string             `bson:"countryCode" json:"countryCode"`
	CountryName string             `bson:"countryName" json:"countryName"`
	TaxType   string             `bson:"taxType" json:"taxType"`
	TaxCode   string             `bson:"taxCode" json:"taxCode"`
	TaxName   string             `bson:"taxName" json:"taxName"`
	TaxRate   float64            `bson:"taxRate" json:"taxRate"`
	ImportRate float64           `bson:"importRate" json:"importRate"`
	Remark    string             `bson:"remark" json:"remark"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}