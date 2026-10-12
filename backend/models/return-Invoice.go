package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReturnInvoice struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	InvoiceNumber    string             `json:"invoiceNumber" bson:"invoiceNumber"`
	ReturnDate       string             `json:"returnDate" bson:"returnDate"`
	ReturnReason     string             `json:"returnReason" bson:"returnReason"`
	ReturnLocationId string             `json:"returnLocationId" bson:"returnLocationId"`
	ReturnCaseNumber string             `json:"returnCaseNumber" bson:"returnCaseNumber"`
	ReturnDetail     string             `json:"returnDatail" bson:"returnDatail"`
	ReturnMethod     string             `json:"returnMethod" bson:"returnMethod"`
	ProcessMethod    string             `json:"processMethod" bson:"processMethod"`
	Refund           bool               `json:"refund" bson:"refund"`
	RefundAmount     float64            `json:"refundAmount" bson:"refundAmount"`
	RefundMethod     string             `json:"refundMethod" bson:"refundMethod"`
	ReturnToVendor   bool               `json:"returnToVendor" bson:"returnToVendor"`
	Status           int                `json:"status" bson:"status"`
	CreatedAt        string             `json:"createdAt" bson:"createdAt"`
	UpdatedAt        string             `json:"updatedAt" bson:"updatedAt"`
}

type ReturnInvoiceItem struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ReturnDataId    string             `json:"returnDataId" bson:"returnDataId"`
	ReturnInvoiceId string             `json:"returnInvoiceId" bson:"returnInvoiceId"`
	ProductId       string             `json:"productId" bson:"productId"`
	Qty             float64            `json:"qty" bson:"qty"`
	Price           float64            `json:"price" bson:"price"`
	TaxType         string             `json:"taxType" bson:"taxType"`
	TaxCode         string             `json:"taxCode" bson:"taxCode"`
	TaxRate         float64            `json:"taxRate" bson:"taxRate"`
	TaxAmount       float64            `json:"taxAmount" bson:"taxAmount"`
}
