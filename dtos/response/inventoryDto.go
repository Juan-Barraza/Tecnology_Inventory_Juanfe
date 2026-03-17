package dtos

import "time"

type InventoryPeriodResponse struct {
	ID          string     `json:"id"`
	PeriodYear  int        `json:"period_year"`
	PeriodMonth int        `json:"period_month"`
	Status      string     `json:"status"`
	CreatedBy   string     `json:"created_by"`
	ClosedAt    *time.Time `json:"closed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type InventoryRecordDetailResponse struct {
	ID               string    `json:"id"`
	PeriodID         string    `json:"period_id"`
	AssetID          string    `json:"asset_id"`
	AssetCode        string    `json:"asset_code"`
	AssetDescription string    `json:"asset_description"`
	CategoryName     string    `json:"category_name"`
	CityName         string    `json:"city_name"`
	AreaName         *string   `json:"area_name"`
	LogicalStatus    string    `json:"logical_status"`
	Confirmed        bool      `json:"confirmed"`
	Deactivated      bool      `json:"deactivated"`
	Notes            *string   `json:"notes"`
	RecordedByName   string    `json:"recorded_by_name"`
	RecordedAt       time.Time `json:"recorded_at"`
}

type AssetInventoryStatusResponse struct {
	AssetID          string     `json:"asset_id"`
	AssetCode        string     `json:"asset_code"`
	AssetDescription string     `json:"asset_description"`
	CategoryName     string     `json:"category_name"`
	CityName         string     `json:"city_name"`
	AreaName         *string    `json:"area_name"`
	RecordID         *string    `json:"record_id"`
	Confirmed        *bool      `json:"confirmed"`
	Deactivated      *bool      `json:"deactivated"`
	Notes            *string    `json:"notes"`
	RecordedAt       *time.Time `json:"recorded_at"`
}

type PeriodProgressResponse struct {
	Total      int     `json:"total"`
	Reviewed   int     `json:"reviewed"`
	Pending    int     `json:"pending"`
	Percentage float64 `json:"percentage"`
}
