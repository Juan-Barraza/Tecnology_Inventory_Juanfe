package dtos

type CreateAssetRequest struct {
	Code              string   `json:"code"`
	Description       string   `json:"description"`
	CategoryID        int      `json:"category_id"`
	AccountingGroupID int      `json:"accounting_group_id"`
	AssetAccountID    int      `json:"asset_account_id"`
	CityID            int      `json:"city_id"`
	AreaID            *int     `json:"area_id"`
	HistoricalCost    *float64 `json:"historical_cost"`
	ActivationDate    string   `json:"activation_date"` // "YYYY-MM-DD"
	PhysicalStatus    string   `json:"physical_status"`
}

type UpdateAssetRequest struct {
	Code              string   `json:"code"`
	Description       *string  `json:"description"`
	CategoryID        *int     `json:"category_id"`
	AccountingGroupID *int     `json:"accounting_group_id"`
	AssetAccountID    *int     `json:"asset_account_id"`
	CityID            *int     `json:"city_id"`
	AreaID            *int     `json:"area_id"`
	HistoricalCost    *float64 `json:"historical_cost"`
	PhysicalStatus    *string  `json:"physical_status"`
	LogicalStatus     *string  `json:"logical_status"`
}

type UpdateAssetStatusRequest struct {
	LogicalStatus  *string `json:"logical_status"`
	PhysicalStatus *string `json:"physical_status"`
	Notes          *string `json:"notes"`
}

type AssetResponse struct {
	ID                  string   `json:"id"`
	Code                string   `json:"code"`
	Description         string   `json:"description"`
	Category            string   `json:"category"`
	AccountCode         int64    `json:"account_code"`
	OpenLedger          *string  `json:"open_ledger"`
	AccountingGroupName string   `json:"accounting_group_name"`
	AccountingGroupCode int64    `json:"accounting_group_code"`
	CityId              int      `json:"city_id"`
	City                string   `json:"city"`
	Area                *string  `json:"area"`
	HistoricalCost      *float64 `json:"historical_cost"`
	ActivationDate      string   `json:"activation_date"`
	LogicalStatus       string   `json:"logical_status"`
	PhysicalStatus      string   `json:"physical_status"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
}

// AssetFilter holds all query params for the list endpoint
// GET /assets?city_id=1&area_id=2&logical_status=active&from=2022-01-01&to=2024-12-31
type AssetFilter struct {
	CityID            *int    `query:"city_id"`
	AreaID            *int    `query:"area_id"`
	CategoryID        *int    `query:"category_id"`
	AccountingGroupID *int    `query:"accounting_group_id"`
	AssetAccountID    *int    `query:"asset_account_id"`
	LogicalStatus     *string `query:"logical_status"`
	PhysicalStatus    *string `query:"physical_status"`
	From              *string `query:"from"`   // activation_date range start
	To                *string `query:"to"`     // activation_date range end
	Search            *string `query:"search"` // code or description ILIKE
	Page              int     `query:"page"`
	Limit             int     `query:"limit"`
}
