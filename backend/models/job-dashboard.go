package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobDashboard struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Number          string             `bson:"number" json:"number"`
	Category        string             `bson:"category" json:"category"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	JobStartDate    string             `bson:"jobStartDate" json:"jobStartDate"`
	JobEndDate      string             `bson:"jobEndDate" json:"jobEndDate"`
	JobTitleRemain  string             `bson:"jobTitleRemain" json:"jobTitleRemain"`
	LocationIds     []string           `bson:"locationIds" json:"locationIds"`
	Tags            []string           `bson:"tags" json:"tags"`
	Labels          []string           `bson:"labels" json:"labels"`
	Attachments     []AttachmentLists  `bson:"attachments" json:"attachments"`
	Comments        []CommentList      `bson:"comments" json:"comments"`
	CheckLists      []CheckList        `bson:"checkLists" json:"checkLists"`
	UserIds         []string           `bson:"userIds" json:"userIds"`
	JobRelationData JobRelationData    `bson:"jobRelationData" json:"jobRelationData"`
	JobValid        bool               `bson:"jobValid" json:"jobValid"`
	SortOrder       int                `bson:"sortOrder" json:"sortOrder"`
	Status          int                `bson:"status" json:"status"`
	CreatedAt       string             `bson:"createdAt,omitempty" json:"createdAt"`
	CreatedBy       string             `bson:"createdBy,omitempty" json:"createdBy"`
	UpdatedAt       string             `bson:"updatedAt,omitempty" json:"updatedAt"`
}

type CheckList struct {
	Name    string `bson:"name" json:"name"`
	Checked bool   `bson:"checked" json:"checked"`
	Status  int    `bson:"status" json:"status"`
}

type CommentList struct {
	UserId   string `bson:"userId" json:"userId"`
	Username string `bson:"username" json:"username"`
	Avatar   string `bson:"avatar" json:"avatar"`
	Comment  string `bson:"comment" json:"comment"`
	Status   int    `bson:"status" json:"status"`
}

type AttachmentLists struct {
	FileName  string `bson:"fileName" json:"fileName"`
	FileType  string `bson:"fileType" json:"fileType"`
	Base64    string `bson:"base64" json:"base64"`
	UpdatedAt string `bson:"updatedAt" json:"updatedAt"`
	Status    int    `bson:"status" json:"status"`
}

type JobRelationData struct {
	ProductCode             string `bson:"productCode" json:"productCode"`
	TypeCode                string `bson:"typeCode" json:"typeCode"`
	DeptCode                string `bson:"deptCode" json:"deptCode"`
	PlaceCode               string `bson:"placeCode" json:"placeCode"`
	VendorCode              string `bson:"vendorCode" json:"vendorCode"`
	StockTakeName           string `bson:"stockTakeName" json:"stockTakeName"`
	InvoiceNumber           string `bson:"invoiceNumber" json:"invoiceNumber"`
	ReturnInvoiceCaseNumber string `bson:"returnInvoiceCaseNumber" json:"returnInvoiceCaseNumber"`
}
