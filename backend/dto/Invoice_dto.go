package dto

import (
	"golang-pos-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoiceDto struct {
	MemberId       string                   `bson:"memberId,omitempty" json:"memberId"`
	LocationId     string                   `bson:"locationId,omitempty" json:"locationId"`
	TotalAmount    float64                  `bson:"totalAmount,omitempty" json:"totalAmount"`
	Discount       float64                  `bson:"discount,omitempty" json:"discount"`
	DiscountType   string                   `bson:"discountType,omitempty" json:"discountType"`
	TaxTotal       float64                  `bson:"taxTotal,omitempty" json:"taxTotal"`
	TaxRefNo       string                   `bson:"taxRefNo,omitempty" json:"taxRefNo"`
	Remark         string                   `bson:"remark,omitempty" json:"remark"`
	InvoiceItems   []models.InvoiceItemList `bson:"invoiceItems,omitempty" json:"invoiceItems"`
	InvoicePayment []models.InvoicePayment  `bson:"invoicePayments,omitempty" json:"invoicePayments"`
}

type InvoicePureBody struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number       string             `bson:"number,omitempty" json:"number"`
	TotalAmount  float64            `bson:"totalAmount,omitempty" json:"totalAmount"`
	Discount     float64            `bson:"discount,omitempty" json:"discount"`
	DiscountType string             `bson:"discountType,omitempty" json:"discountType"`
	TaxTotal     float64            `bson:"taxTotal,omitempty" json:"taxTotal"`
	TaxRefNo     string             `bson:"taxRefNo,omitempty" json:"taxRefNo"`
	Remark       string             `bson:"remark,omitempty" json:"remark"`
	Status       int                `bson:"status" json:"status"`
	CreatedAt    time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
	PlaceCode    string             `bson:"placeCode" json:"placeCode"`
	PlaceName    string             `bson:"placeName" json:"placeName"`
}

type InvoiceBody struct {
	ID              primitive.ObjectID      `bson:"_id,omitempty" json:"id"`
	Number          string                  `bson:"number,omitempty" json:"number"`
	MemberId        string                  `bson:"memberId,omitempty" json:"memberId"`
	LocationId      string                  `bson:"locationId,omitempty" json:"locationId"`
	TotalAmount     float64                 `bson:"totalAmount,omitempty" json:"totalAmount"`
	Discount        float64                 `bson:"discount,omitempty" json:"discount"`
	DiscountType    string                  `bson:"discountType,omitempty" json:"discountType"`
	TaxTotal        float64                 `bson:"taxTotal,omitempty" json:"taxTotal"`
	TaxRefNo        string                  `bson:"taxRefNo,omitempty" json:"taxRefNo"`
	Remark          string                  `bson:"remark,omitempty" json:"remark"`
	Status          int                     `bson:"status" json:"status"`
	CreatedAt       time.Time               `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt       time.Time               `bson:"updatedAt,omitempty" json:"updatedAt"`
	InvoiceItems    []InvoiceItemListBody   `bson:"invoiceItems,omitempty" json:"invoiceItems"`
	InvoicePayments []models.InvoicePayment `bson:"invoicePayments,omitempty" json:"invoicePayments"`
}

type InvoiceItemListBody struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	InvoiceId    string              `bson:"invoiceId" json:"invoiceId"`
	ProductId    string              `bson:"productId" json:"productId"`
	Qty          float64             `bson:"qty" json:"qty"`
	Price        float64             `bson:"price" json:"price"`
	Discount     float64             `bson:"discount" json:"discount"`
	DiscountType string              `bson:"discountType" json:"discountType"`
	TaxType      string              `bson:"taxType" json:"taxType"`
	TaxCode      string              `bson:"taxCode" json:"taxCode"`
	TaxRate      float64             `bson:"taxRate" json:"taxRate"`
	TaxAmount    float64             `bson:"taxAmount" json:"taxAmount"`
	Product      models.ProductLists `bson:"product" json:"product"`
}

type InvoiceListBody struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number       string             `bson:"number,omitempty" json:"number"`
	MemberId     string             `bson:"memberId,omitempty" json:"memberId"`
	LocationId   string             `bson:"locationId,omitempty" json:"locationId"`
	TotalAmount  float64            `bson:"totalAmount,omitempty" json:"totalAmount"`
	Discount     bool               `bson:"discount,omitempty" json:"discount"`
	DiscountType string             `bson:"discountType,omitempty" json:"discountType"`
	TaxTotal     float64            `bson:"taxTotal,omitempty" json:"taxTotal"`
	TaxRefNo     string             `bson:"taxRefNo,omitempty" json:"taxRefNo"`
	Remark       string             `bson:"remark,omitempty" json:"remark"`
	Status       int                `bson:"status" json:"status"`
	CreatedAt    time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
	Location     models.Locations   `bson:"location" json:"location"`
}

type ListInvoiceRequestDto struct {
	Page        int      `bson:"page" json:"page"`
	Limit       int      `bson:"limit" json:"limit"`
	LocationIds []string `bson:"locationId" json:"locationId"`
	Date        []string `bson:"date" json:"date"`
}

type ListInvoices struct {
	Lists []InvoiceListBody `bson:"lists" json:"lists"`
	Total int               `bson:"total" json:"total"`
	Page  int               `bson:"page" json:"page"`
	Limit int               `bson:"limit" json:"limit"`
}

type DashboardReqFilterDto struct {
	ProductCode    string   `json:"productCode"`
	ProductName    string   `json:"productName"`
	TypeIds        []string `json:"typeIds"`
	PlaceIds       []string `json:"placeIds"`
	DeptIds        []string `json:"deptIds"`
	SalesDateRange []string `json:"salesDateRange"`
}

type InvoiceDashboardReqDto struct {
	DataTypeValue string                `json:"dataTypeValue"`
	DateTypeValue string                `json:"dateTypeValue"`
	ValueField    string                `json:"valueField"`
	Filter        DashboardReqFilterDto `json:"filter"`
}

type InvRecordYearMonthResult struct {
	Year     int     `bson:"year" json:"year"`
	Month    string  `bson:"month" json:"month"`
	MonthNum int     `bson:"monthNum" json:"monthNum"`
	Qtys     float64 `bson:"qtys" json:"qtys"`
	Costs    float64 `bson:"costs" json:"costs"`
}
