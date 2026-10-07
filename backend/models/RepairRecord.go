package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RepairRecords struct {
	ID   	  		      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssetId               string             `bson:"assetId" json:"assetId"`
	RepairReason  	      string             `bson:"repairReason" json:"repairReason"`
	MaintenancePeriod     bool               `bson:"maintenancePeriod" json:"maintenancePeriod"`
	MaintenanceName		  string             `bson:"maintenanceName" json:"maintenanceName"`
	MaintenanceDate       string             `bson:"maintenanceDate" json:"maintenanceDate"` //MaintenanceDate
	MaintenanceFinishDate string         `bson:"maintenanceFinishDate" json:"maintenanceFinishDate"`
	RepairInvoiceDate	  string         `bson:"repairInvoiceDate" json:"repairInvoiceDate"`
	RepairInvoiceNo 	  string             `bson:"repairInvoiceNo" json:"repairInvoiceNo"`
	RepairAmount          float64            `bson:"repairAmount" json:"repairAmount"` //RepairAmount
	Remark				  string             `bson:"remark" json:"remark"`
	Status                int                `bson:"status" json:"status"`
	CreatedAt             time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt             time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}