package dto

import (
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolesPageDto struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
}

type RolesList struct {
	Lists []models.SysRoles `json:"lists"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

type MenuItemPermissionBody struct {
	MenuIds []any  `json:"menuIds"`
	ID      string `json:"id"`
}

type RoleIdsBody struct {
	RoleIds []string `json:"roleIds"`
}

type ListAllSysRoleDto struct {
	Data []models.SysRoles `json:"data"`
}

type ListMenus struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	MainId            string             `bson:"mainId" json:"mainId"`
	Name              string             `bson:"name" json:"name"`
	Icon              string             `bson:"icon" json:"icon"`
	Path              string             `bson:"path" json:"path"`
	Sort              int                `bson:"sort" json:"sort"`
	Type              string             `bson:"type" json:"type"`
	ExcelFunctionCode string             `bson:"excelFunctionCode" json:"excelFunctionCode"`
	ExcelFunctionName string             `bson:"excelFunctionName" json:"excelFunctionName"`
	Status            int                `bson:"status" json:"status"`
	CreatedAt         time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
	Read              bool               `bson:"read" json:"read"`
	Write             bool               `bson:"write" json:"write"`
	Delete            bool               `bson:"delete" json:"delete"`
	Upload            bool               `bson:"upload" json:"upload"`
	Update            bool               `bson:"update" json:"update"`
}

type SysRolesWithMenus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	Name      string             `bson:"name" json:"name"`
	Code      string             `bson:"code" json:"code"`
	Remark    string             `bson:"remark" json:"remark"`
	MenuIds   []any              `bson:"menuIds" json:"menuIds"`
	Read      bool               `bson:"read" json:"read"`
	Write     bool               `bson:"write" json:"write"`
	Delete    bool               `bson:"delete" json:"delete"`
	Upload    bool               `bson:"upload" json:"upload"`
	Update    bool               `bson:"update" json:"update"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
	MenuLists []ListMenus        `bson:"menuLists" json:"menuLists"`
}
