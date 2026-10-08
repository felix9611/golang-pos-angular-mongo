package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductLists struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProductCode            string             `bson:"productCode" json:"productCode"`
	ProductName            string             `bson:"productName" json:"productName"`
	ItemCode               string             `bson:"itemCode" json:"itemCode"`
	BrandCode              string             `bson:"brandCode" json:"brandCode"`
	BrandName              string             `bson:"brandName" json:"brandName"`
	TypeID                 string             `bson:"typeId" json:"typeId"`
	DeptID                 string             `bson:"deptId" json:"deptId"`
	Unit                   string             `bson:"unit" json:"unit"`
	CostPrice              float64            `bson:"costPrice" json:"costPrice"`
	RetailPrice            float64            `bson:"retailPrice" json:"retailPrice"`
	Description            string             `bson:"description" json:"description"`
	TaxType                string             `bson:"taxType" json:"taxType"`
	Remark                 string             `bson:"remark" json:"remark"`
	CreatedAt              time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt              time.Time          `bson:"updated_at" json:"updatedAt"`
	Status                 int                `bson:"status" json:"status"`
	UploadProductListFiles []ProductListFiles `bson:"uploadProductListFiles" json:"uploadProductListFiles"`
}

type ProductListFiles struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID string             `bson:"assetId" json:"assetId"`
	FileName  string             `bson:"fileName" json:"fileName"`
	FileType  string             `bson:"fileType" json:"fileType"`
	Base64    string             `bson:"base64" json:"base64"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt string             `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt string             `bson:"updatedAt,omitempty" json:"updatedAt"`
}
