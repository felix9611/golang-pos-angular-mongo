package dto

import "time"

type RepairRecordPageReqDTO struct {
	Page      int64    `json:"page"`
	Limit     int64    `json:"limit"`
	DateRange []string `json:"dateRange"`
	AssetCode string   `json:"assetCode"`
	DeptIds   []string `json:"deptIds"`
	TypeIds   []string `json:"typeIds"`
	PlaceIds  []string `json:"placeIds"`
}

type RepairRecordPureList struct {
	AssetCode             string    `json:"assetCode"`
	AssetName             string    `json:"assetName"`
	RepairReason          string    `json:"repairReason"`
	MaintenanceReriod     string    `json:"maintenanceReriod"`
	MaintenanceDate       string    `json:"maintenanceDate"`
	MaintenanceFinishDate string    `json:"maintenanceFinishDate"`
	RepairInvoiceDate     string    `json:"repairInvoiceDate"`
	MaintenanceName       string    `json:"maintenanceName"`
	RepairAmount          float64   `json:"repairAmount"`
	Remark                string    `json:"remark"`
	CreatedAt             time.Time `json:"createdAt"`
}
