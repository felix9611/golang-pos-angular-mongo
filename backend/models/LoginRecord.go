package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LoginRecords struct {
	ID        	primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  	string             `bson:"username" json:"username"`
	IpAddress 	string             `bson:"ipAddress" json:"ipAddress"`
	LoginStatus string            	`bson:"loginStatus" json:"loginStatus"`
	LoginTime 	time.Time          `bson:"loginTime" json:"loginTime"`
}