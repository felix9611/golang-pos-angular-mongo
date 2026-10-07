package dto

import (
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListBudgetRecordsDto struct {
	Name    string      `json:"name"`
	Page    int64       `json:"page"`
	DeptID  []string    `json:"deptId"`
	PlaceId []string    `json:"placeId"`
	Date    []time.Time `json:"date"`
	Limit   int64       `json:"limit"`
}

type DateRange struct {
	From primitive.DateTime `json:"from"`
	To   primitive.DateTime `json:"to"`
}

type BudgetSummary struct {
	BudgetAmount float64 `bson:"budgetAmount" json:"budgetAmount"`
	YearMonth    string  `bson:"yearMonth" json:"yearMonth"`
	Year         int     `bson:"year" json:"year"`
	Month        int     `bson:"month" json:"month"`
}

type BudgetList struct {
	Lists []models.Budgets `json:"lists"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

type BudgetPureList struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DeptID       primitive.ObjectID `bson:"deptId" json:"deptId"`
	PlaceId      primitive.ObjectID `bson:"placeId" json:"placeId"`
	BudgetNo     string             `bson:"budgetNo" json:"budgetNo"`
	BudgetName   string             `bson:"budgetName" json:"budgetName"`
	Year         string             `bson:"year" json:"year"`
	Month        string             `bson:"month" json:"month"`
	BudgetAmount float64            `bson:"budgetAmount" json:"budgetAmount"`
	BudgetFrom   primitive.DateTime `bson:"budgetFrom" json:"budgetFrom"`
	BudgetTo     primitive.DateTime `bson:"budgetTo" json:"budgetTo"`
	BudgetStatus string             `bson:"budgetStatus" json:"budgetStatus"`
	Remark       string             `bson:"remark" json:"remark"`
	Status       int                `bson:"status" json:"status"`
	CreatedAt    time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
	PlaceName    string             `bson:"placeName" json:"placeName"`
	PlaceCode    string             `bson:"placeCode" json:"placeCode"`
	DeptName     string             `bson:"deptName" json:"deptName"`
	DeptCode     string             `bson:"deptCode" json:"deptCode"`
}
