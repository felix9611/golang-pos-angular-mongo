package dto

import (
	"golang-pos-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListProductDetail struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	ProductCode string              `bson:"productCode"`
	ProductName string              `bson:"productName"`
	ItemCode    string              `bson:"itemCode"`
	BrandCode   string              `bson:"brandCode"`
	BrandName   string              `bson:"brandName"`
	TypeID      string              `bson:"typeId"`
	DeptID      string              `bson:"deptId"`
	Unit        string              `bson:"unit"`
	CostPrice   float64             `bson:"costPrice"`
	RetailPrice float64             `bson:"retailPrice"`
	Description string              `bson:"description"`
	TaxType     string              `bson:"taxType"`
	Remark      string              `bson:"remark"`
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
	Status      int                 `bson:"status"`
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
