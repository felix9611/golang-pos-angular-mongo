package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SysRoles struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	Name        string             `bson:"name" json:"name"`
	Code 	  string             `bson:"code" json:"code"`
	Remark	  string             `bson:"remark" json:"remark"`
	MenuIds []any	`bson:"menuIds" json:"menuIds"`
	Read  bool               `bson:"read" json:"read"`
	Write bool               `bson:"write" json:"write"`
	Delete bool             `bson:"delete" json:"delete"`
	Upload bool             `bson:"upload" json:"upload"`
	Update bool             `bson:"update" json:"update"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}