package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SysUsers struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"  json:"id"`
	Username     string               `bson:"username" json:"username"`
	Password     string               `bson:"password" json:"password"`
	Email        string               `bson:"email" json:"email"`
	AvatarBase64 string               `bson:"avatarBase64" json:"avatarBase64"`
	DeptIds      []primitive.ObjectID `bson:"DeptIds,omitempty" json:"deptIds"`
	Roles        []primitive.ObjectID `bson:"roles,omitempty" json:"roles"`
	SuperAdmin   bool                 `bson:"superAdmin" json:"superAdmin"`
	PlaceIds     []primitive.ObjectID `bson:"placeIds" json:"placeIds"`
	Status       int                  `bson:"status" json:"status"`
	CreatedAt    time.Time            `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt    time.Time            `bson:"updatedAt,omitempty" json:"updatedAt"`
}
