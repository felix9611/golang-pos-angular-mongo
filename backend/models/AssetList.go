package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	//	"time"
)

type AssetLists struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssetCode              string             `bson:"assetCode" json:"assetCode"`
	AssetName              string             `bson:"assetName" json:"assetName"`
	Unit                   string             `bson:"unit" json:"unit"`
	TypeID                 string             `bson:"typeId" json:"typeId"`
	PlaceId                string             `bson:"placeId" json:"placeId"`
	DeptId                 string             `bson:"deptId" json:"deptId"`
	PurchaseDate           string             `bson:"purchaseDate" json:"purchaseDate"`
	Description            string             `bson:"description" json:"description"`
	Sponsor                bool               `bson:"sponsor" json:"sponsor"`
	SponsorName            string             `bson:"sponsorName" json:"sponsorName"`
	Cost                   float64            `bson:"cost" json:"cost"`
	SerialNumber           string             `bson:"serialNo" json:"serialNo"`
	InvoiceNo              string             `bson:"invoiceNo" json:"invoiceNo"`
	InvoiceDate            string             `bson:"invoiceDate" json:"invoiceDate"`
	InvoiceRemark          string             `bson:"invoiceRemark" json:"invoiceRemark"` // I
	VendorId               string             `bson:"vendorId" json:"vendorId"`
	Remark                 string             `bson:"remark" json:"remark"`
	TaxInfoId              string             `bson:"taxInfoId" json:"taxInfoId"`
	TaxCountryCode         string             `bson:"taxCountryCode" json:"taxCountryCode"`
	TaxCode                string             `bson:"taxCode" json:"taxCode"`
	TaxRate                float64            `bson:"taxRate" json:"taxRate"`
	IncludeTax             bool               `bson:"includeTax" json:"includeTax"`
	AfterBeforeTax         float64            `bson:"afterBeforeTax" json:"afterBeforeTax"`
	AccountCode            string             `bson:"accountCode" json:"accountCode"`
	AccountName            string             `bson:"accountName" json:"accountName"`
	BrandCode              string             `bson:"brandCode" json:"brandCode"`
	BrandName              string             `bson:"brandName" json:"brandName"`
	ChequeNo               string             `bson:"chequeNo" json:"chequeNo"`
	MaintenancePeriodStart string             `bson:"maintenancePeriodStart" json:"maintenancePeriodStart"`
	MaintenancePeriodEnd   string             `bson:"maintenancePeriodEnd" json:"maintenancePeriodEnd"`
	VoucherNo              string             `bson:"voucherNo" json:"voucherNo"`
	VoucherUsedDate        string             `bson:"voucherUsedDate" json:"voucherUsedDate"`
	StaffName              string             `bson:"staffName" json:"staffName"`
	Status                 int                `bson:"status" json:"status"`
	CreatedAt              string             `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt              string             `bson:"updatedAt,omitempty" json:"updatedAt"`
	UploadAssetListFiles   []AssetListFiles   `bson:"uploadAssetListFiles" json:"uploadAssetListFiles"`
}

type AssetListFiles struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssetId   string             `bson:"assetId" json:"assetId"`
	FileName  string             `bson:"fileName" json:"fileName"`
	FileType  string             `bson:"fileType" json:"fileType"`
	Base64    string             `bson:"base64" json:"base64"`
	Status    int                `bson:"status" json:"status"`
	CreatedAt string             `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt string             `bson:"updatedAt,omitempty" json:"updatedAt"`
}
