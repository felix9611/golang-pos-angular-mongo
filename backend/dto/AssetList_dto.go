package dto

import (
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetListsPureDetails struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssetCode              string             `bson:"assetCode" json:"assetCode"`
	AssetName              string             `bson:"assetName" json:"assetName"`
	Unit                   string             `bson:"unit" json:"unit"`
	TypeID                 string             `bson:"typeId" json:"typeId"`
	PlaceId                string             `bson:"placeId" json:"placeId"`
	DeptId                 string             `bson:"deptId" json:"deptId"`
	PurchaseDate           string             `bson:"purchaseDate" json:"purchaseDate"`
	Description            string             `bson:"description" json:"description"`
	Sponsor                string             `bson:"sponsor" json:"sponsor"`
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
	IncludeTax             string             `bson:"includeTax" json:"includeTax"`
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
	PlaceCode              string             `bson:"placeCode" json:"placeCode"`
	PlaceName              string             `bson:"placeName" json:"placeName"`
	TypeCode               string             `bson:"typeCode" json:"typeCode"`
	TypeName               string             `bson:"typeName" json:"typeName"`
	DeptCode               string             `bson:"deptCode" json:"deptCode"`
	DeptName               string             `bson:"deptName" json:"deptName"`
}

type AssetListsDetails struct {
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
	Location               models.Locations   `bson:"location" json:"location"`
	Department             models.Department  `bson:"department" json:"department"`
	AssetType              models.AssetTypes  `bson:"assetType" json:"assetType"`
}

type ListAssetReqDto struct {
	Page          int      `json:"page"`
	Limit         int      `json:"limit"`
	AssetCode     string   `json:"assetCode"`
	AssetName     string   `json:"assetName"`
	TypeIds       []string `json:"typeIds"`
	PlaceIds      []string `json:"placeIds"`
	DeptIds       []string `json:"deptIds"`
	PurchaseDates []string `json:"purchaseDates"`
}

type ListRecordReqDto struct {
	Page      int      `json:"page"`
	Limit     int      `json:"limit"`
	AssetCode string   `json:"assetCode"`
	DateRange []string `json:"dateRange"`
}

type AssetListList struct {
	Lists []AssetListsDetails `json:"lists"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}

type DashboardReqFilterDto struct {
	TypeIds       []string    `json:"typeIds"`
	PlaceIds      []string    `json:"placeIds"`
	DeptIds       []string    `json:"deptIds"`
	PurchaseDates []time.Time `json:"purchaseDates"`
}

type DashboardReqDto struct {
	DataType      bool                   `json:"dataType"`
	DataTypeValue string                 `json:"dataTypeValue"`
	DateType      bool                   `json:"dateType"`
	DateTypeValue string                 `json:"dateTypeValue"`
	ValueField    string                 `json:"valueField"`
	Filter        *DashboardReqFilterDto `json:"filter"`
}
