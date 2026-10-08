package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductLists struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty"`
	ProductCode            string             `bson:"productCode"`
	ProductName            string             `bson:"productName"`
	ItemCode               string             `bson:"itemCode"`
	BrandCode              string             `bson:"brandCode"`
	BrandName              string             `bson:"brandName"`
	TypeID                 string             `bson:"typeId"`
	DeptID                 string             `bson:"deptId"`
	Unit                   string             `bson:"unit"`
	CostPrice              float64            `bson:"costPrice"`
	RetailPrice            float64            `bson:"retailPrice"`
	Description            string             `bson:"description"`
	TaxType                string             `bson:"taxType"`
	Remark                 string             `bson:"remark"`
	CreatedAt              time.Time          `bson:"created_at"`
	UpdatedAt              time.Time          `bson:"updated_at"`
	Status                 int                `bson:"status"`
	UploadProductListFiles []ProductListFiles `bson:"uploadProductListFiles"`
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
