package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Number       string             `bson:"number" json:"number"`
	MemberId     string             `bson:"memberId" json:"memberId"`
	LocationId   string             `bson:"locationId" json:"locationId"`
	TotalAmount  float64            `bson:"totalAmount" json:"totalAmount"`
	Discount     float64            `bson:"discount" json:"discount"`
	DiscountType string             `bson:"discountType" json:"discountType"`
	TaxTotal     float64            `bson:"taxTotal" json:"taxTotal"`
	TaxRefNo     string             `bson:"taxRefNo" json:"taxRefNo"`
	Remark       string             `bson:"remark" json:"remark"`
	Status       int                `bson:"status" json:"status"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type InvoiceItemList struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	InvoiceId    string             `bson:"invoiceId" json:"invoiceId"`
	ProductId    string             `bson:"productId" json:"productId"`
	Qty          float64            `bson:"qty" json:"qty"`
	Price        float64            `bson:"price" json:"price"`
	Discount     float64            `bson:"discount" json:"discount"`
	DiscountType string             `bson:"discountType" json:"discountType"`
	TaxType      string             `bson:"taxType" json:"taxType"`
	TaxCode      string             `bson:"taxCode" json:"taxCode"`
	TaxRate      float64            `bson:"taxRate" json:"taxRate"`
	TaxAmount    float64            `bson:"taxAmount" json:"taxAmount"`
}

type InvoicePayment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	InvoiceId   string             `bson:"invoiceId" json:"invoiceId"`
	Method      string             `bson:"method" json:"method"`
	Amount      float64            `bson:"amount" json:"amount"`
	PaymentTime string             `bson:"paymentTime" json:"paymentTime"`
}
