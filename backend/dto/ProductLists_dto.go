package dto

import (
	"golang-pos-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListProductDetailPure struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductCode string             `bson:"productCode" json:"productCode"`
	ProductName string             `bson:"productName" json:"productName"`
	ItemCode    string             `bson:"itemCode" json:"itemCode"`
	BrandCode   string             `bson:"brandCode" json:"brandCode"`
	BrandName   string             `bson:"brandName" 	json:"brandName"`
	TypeID      string             `bson:"typeId"	json:"typeId"`
	DeptID      string             `bson:"deptId" 	json:"deptId"`
	Unit        string             `bson:"unit" json:"unit"`
	CostPrice   float64            `bson:"costPrice" json:"costPrice"`
	RetailPrice float64            `bson:"retailPrice" json:"retailPrice"`
	Description string             `bson:"description" json:"description"`
	TaxType     string             `bson:"taxType" json:"taxType"`
	Remark      string             `bson:"remark" json:"remark"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
	Status      int                `bson:"status" json:"status"`
	DeptCode    string             `bson:"deptCode" json:"deptCode"`
	DeptName    string             `bson:"deptName" json:"deptName"`
	TypeCode    string             `bson:"typeCode" json:"typeCode"`
	TypeName    string             `bson:"typeName" json:"typeName"`
}

type ListProductDetail struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ProductCode string              `bson:"productCode" json:"productCode"`
	ProductName string              `bson:"productName" json:"productName"`
	ItemCode    string              `bson:"itemCode" json:"itemCode"`
	BrandCode   string              `bson:"brandCode" json:"brandCode"`
	BrandName   string              `bson:"brandName" 	json:"brandName"`
	TypeID      string              `bson:"typeId"	json:"typeId"`
	DeptID      string              `bson:"deptId" 	json:"deptId"`
	Unit        string              `bson:"unit" json:"unit"`
	CostPrice   float64             `bson:"costPrice" json:"costPrice"`
	RetailPrice float64             `bson:"retailPrice" json:"retailPrice"`
	Description string              `bson:"description" json:"description"`
	TaxType     string              `bson:"taxType" json:"taxType"`
	Remark      string              `bson:"remark" json:"remark"`
	CreatedAt   time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updatedAt"`
	Status      int                 `bson:"status" json:"status"`
	Department  models.Department   `json:"department" bson:"department"`
	ProductType models.ProductTypes `json:"productType" bson:"productType"`
}

type ListProductReqDto struct {
	Page        int      `json:"page"`
	Limit       int      `json:"limit"`
	ProductCode string   `json:"assetCode"`
	ProductName string   `json:"assetName"`
	TypeIds     []string `json:"typeIds"`
	DeptIds     []string `json:"deptIds"`
}

type ProductListList struct {
	Lists []ListProductDetail `json:"lists"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}
