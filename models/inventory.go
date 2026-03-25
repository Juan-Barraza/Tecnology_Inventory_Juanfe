package models

import "time"

type PeriodStatus string

const (
	PeriodStatusOpen   PeriodStatus = "open"
	PeriodStatusClosed PeriodStatus = "closed"
)

type InventoryPeriod struct {
	ID          string
	PeriodYear  int
	PeriodMonth int
	PeriodDay   int
	Status      PeriodStatus
	CreatedBy   string
	ClosedAt    *time.Time
	CreatedAt   time.Time
}

type InventoryRecord struct {
	ID          string
	PeriodID    string
	AssetID     string
	Confirmed   bool
	Deactivated bool
	Notes       *string
	HasLabel    bool
	RecordedBy  string
	RecordedAt  time.Time
}

// InventoryRecordDetail incluye datos del asset para mostrarlo en el frontend
type InventoryRecordDetail struct {
	InventoryRecord
	AssetCode        string
	AssetDescription string
	CategoryName     string
	CityName         string
	AreaName         *string
	LogicalStatus    string
	RecordedByName   string
}

type AssetInventoryStatus struct {
	AssetID          string
	AssetCode        string
	AssetDescription string
	CategoryName     string
	CityName         string
	AreaName         *string
	RecordID         *string
	Confirmed        *bool
	Deactivated      *bool
	HasLabel         *bool
	ActivationDate   time.Time
	Notes            *string
	RecordedAt       *time.Time
}
