package models

import "time"

// LogicalStatus maps to logical_status_enum in PostgreSQL
type LogicalStatus string

const (
	LogicalStatusActive     LogicalStatus = "active"
	LogicalStatusInactive   LogicalStatus = "inactive"
	LogicalStatusWrittenOff LogicalStatus = "written_off"
)

// PhysicalStatus maps to physical_status_enum in PostgreSQL
type PhysicalStatus string

const (
	PhysicalStatusOptimal      PhysicalStatus = "optimal"
	PhysicalStatusGood         PhysicalStatus = "good"
	PhysicalStatusFair         PhysicalStatus = "fair"
	PhysicalStatusDeteriorated PhysicalStatus = "deteriorated"
	PhysicalStatusOutOfService PhysicalStatus = "out_of_service"
)

type Asset struct {
	ID             string
	Code           string
	Description    string
	CategoryID     int
	AssetAccountID int
	CityID         int
	AreaID         *int
	HistoricalCost *float64
	ActivationDate time.Time
	LogicalStatus  LogicalStatus
	PhysicalStatus PhysicalStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// AssetDetail is used for list/detail queries with joined fields
// so handlers never need to do extra lookups
type AssetDetail struct {
	Asset
	CategoryName        string
	AccountingGroupName string
	AccountingGroupCode int64
	AccountCode         int64
	OpenLedger          *string
	CityName            string
	AreaName            *string
}

// AssetAccount is the specific sub-account (account_code is unique).
// Each asset points to one AssetAccount, not directly to AccountingGroup.
type AssetAccount struct {
	ID                int
	AccountingGroupID int
	AccountCode       int64
	OpenLedger        *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type AssetCategory struct {
	ID   int
	Name string
}

type StatusHistory struct {
	ID             string
	AssetID        string
	PreviousStatus *LogicalStatus
	NewStatus      LogicalStatus
	Notes          *string
	RecordedBy     string
	CreatedAt      time.Time
}

type StatusHistoryDetail struct {
	StatusHistory
	RecordedByName string
}
