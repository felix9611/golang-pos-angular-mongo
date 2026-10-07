package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Vendors struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VendorCode string             `bson:"vendorCode" json:"vendorCode"`
	VendorName string             `bson:"vendorName" json:"vendorName"`
	VendorOtherName string             `bson:"vendorOtherName" json:"vendorOtherName"` // v
	Type      string             `bson:"type" json:"type"`
	Email     string             `bson:"email" json:"email"`
	Phone     string             `bson:"phone" json:"phone"`
	Address   string             `bson:"address" json:"address"`
	Fax       string             `bson:"fax" json:"fax"`
	ContactPerson string          `bson:"contactPerson" json:"contactPerson"`
	Website   string             `bson:"website" json:"website"`
	Remark    string             `bson:"remark" json:"remark"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}
