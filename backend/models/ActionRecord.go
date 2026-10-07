package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ActionRecords struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActionName      string             `bson:"actionName" json:"actionName"`
	ActionMethod    string             `bson:"actionMethod" json:"actionMethod"`
	ActionFrom      string             `bson:"actionFrom" json:"actionFrom"`
	ActionData      any                `bson:"actionData" json:"actionData"`
	ActionSuccess   string             `bson:"actionSuccess" json:"actionSuccess"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
}