package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRecords struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username"`
	IpAddress   string             `bson:"ipAddress" json:"ipAddress"`
	Ray         string             `bson:"ray" json:"ray"`
	LoginStatus string             `bson:"loginStatus" json:"loginStatus"`
	LoginTime   time.Time          `bson:"loginTime" json:"loginTime"`
}
