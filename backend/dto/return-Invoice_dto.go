package dto

import (
	"golang-fixedasset-mongo-backend/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateReturnInvoiceItemDto struct {
	ReturnInvoiceId string  `json:"returnInvoiceId" bson:"returnInvoiceId"`
	ProductId       string  `json:"productId" bson:"productId"`
	Qty             float64 `json:"qty" bson:"qty"`
	Price           float64 `json:"price" bson:"price"`
	TaxType         string  `json:"taxType" bson:"taxType"`
	TaxCode         string  `json:"taxCode" bson:"taxCode"`
	TaxRate         float64 `json:"taxRate" bson:"taxRate"`
	TaxAmount       float64 `json:"taxAmount" bson:"taxAmount"`
}

type CreateReturnInvoiceDto struct {
	ReturnInvoiceNo  string  `json:"returnInvoiceNo" bson:"returnInvoiceNo"`
	ReturnReason     string  `json:"returnReason" bson:"returnReason"`
	ReturnLocationId string  `json:"returnLocationId" bson:"returnLocationId"`
	ReturnDatail     string  `json:"returnDatail" bson:"returnDatail"`
	ReturnMethod     string  `json:"returnMethod" bson:"returnMethod"`
	ProcessMethod    string  `json:"processMethod" bson:"processMethod"`
	Refund           bool    `json:"refund" bson:"refund"`
	RefundAmount     float64 `json:"refundAmount" bson:"refundAmount"`
	RefundMethod     string  `json:"refundMethod" bson:"refundMethod"`
	ReturnToVendor   bool    `json:"returnToVendor" bson:"returnToVendor"`
	ReturnItems      []CreateReturnInvoiceItemDto
}
type ReturnInvoiceItemDetail struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ReturnDataId    string              `json:"returnDataId" bson:"returnDataId"`
	ReturnInvoiceId string              `json:"returnInvoiceId" bson:"returnInvoiceId"`
	ProductId       string              `json:"productId" bson:"productId"`
	Qty             float64             `json:"qty" bson:"qty"`
	Price           float64             `json:"price" bson:"price"`
	TaxType         string              `json:"taxType" bson:"taxType"`
	TaxCode         string              `json:"taxCode" bson:"taxCode"`
	TaxRate         float64             `json:"taxRate" bson:"taxRate"`
	TaxAmount       float64             `json:"taxAmount" bson:"taxAmount"`
	Product         models.ProductLists `bson:"product" json:"product"`
}

type ReturnInvoiceDetail struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	InvoiceNumber    string             `json:"invoiceNumber" bson:"invoiceNumber"`
	ReturnDate       string             `json:"returnDate" bson:"returnDate"`
	ReturnReason     string             `json:"returnReason" bson:"returnReason"`
	ReturnLocationId string             `json:"returnLocationId" bson:"returnLocationId"`
	ReturnCaseNumber string             `json:"returnCaseNumber" bson:"returnCaseNumber"`
	ReturnDatail     string             `json:"returnDatail" bson:"returnDatail"`
	ReturnMethod     string             `json:"returnMethod" bson:"returnMethod"`
	ProcessMethod    string             `json:"processMethod" bson:"processMethod"`
	Refund           bool               `json:"refund" bson:"refund"`
	RefundAmount     float64            `json:"refundAmount" bson:"refundAmount"`
	RefundMethod     string             `json:"refundMethod" bson:"refundMethod"`
	ReturnToVendor   bool               `json:"returnToVendor" bson:"returnToVendor"`
	Invoice          models.Invoice     `json:"invoice" bson:"invoice"`
	//	Member           models.MemberLists `json:"member" bson:"member"`
	Location  models.Locations `json:"location" bson:"location"`
	CreatedAt string           `json:"createdAt" bson:"createdAt"`
	UpdatedAt string           `json:"updatedAt" bson:"updatedAt"`
}

type ReturnInvoiceDetailResponse struct {
	ID               primitive.ObjectID        `bson:"_id,omitempty" json:"id"`
	InvoiceNumber    string                    `json:"invoiceNumber" bson:"invoiceNumber"`
	ReturnDate       string                    `json:"returnDate" bson:"returnDate"`
	ReturnReason     string                    `json:"returnReason" bson:"returnReason"`
	ReturnLocationId string                    `json:"returnLocationId" bson:"returnLocationId"`
	ReturnCaseNumber string                    `json:"returnCaseNumber" bson:"returnCaseNumber"`
	ReturnDatail     string                    `json:"returnDatail" bson:"returnDatail"`
	ReturnMethod     string                    `json:"returnMethod" bson:"returnMethod"`
	ProcessMethod    string                    `json:"processMethod" bson:"processMethod"`
	Refund           bool                      `json:"refund" bson:"refund"`
	RefundAmount     float64                   `json:"refundAmount" bson:"refundAmount"`
	RefundMethod     string                    `json:"refundMethod" bson:"refundMethod"`
	ReturnToVendor   bool                      `json:"returnToVendor" bson:"returnToVendor"`
	Invoice          models.Invoice            `json:"invoice" bson:"invoice"`
	Member           models.MemberLists        `json:"member" bson:"member"`
	ReturnItems      []ReturnInvoiceItemDetail `json:"returnItems" bson:"returnItems"`
}

type InvoiceListRequestDto struct {
	Page          int      `json:"page"`
	Limit         int      `json:"limit"`
	CaseNumber    string   `json:"caseNumber"`
	InvoiceNumber string   `json:"invoiceNumber"`
	DateRange     []string `json:"dateRange"`
}

type ReturnInvoiceList struct {
	Lists []ReturnInvoiceDetail `json:"lists"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}
