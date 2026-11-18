package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Login struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	IpAddress string `form:"ipAddress" json:"ipAddress" binding:"required"`
}

type AuthUser struct {
	Username   string               `json:"username"`
	RoleIds    []primitive.ObjectID `json:"roleIds"`
	PlaceIds   []primitive.ObjectID `json:"placeIds"`
	DeptIds    []primitive.ObjectID `json:"deptIds"`
	SuperAdmin bool                 `json:"superAdmin"`
	Manager    bool                 `json:"manager"`
	ID         string               `json:"id"`
}
