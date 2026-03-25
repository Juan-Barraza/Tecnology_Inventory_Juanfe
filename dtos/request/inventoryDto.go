package dtos

type CreatePeriodRequest struct {
	PeriodYear  int `json:"period_year"`
	PeriodMonth int `json:"period_month"`
	PeriodDay   int `json:"period_day"`
}

type RecordAssetRequest struct {
	PeriodID    string  `json:"period_id"`
	AssetID     string  `json:"asset_id"`
	Confirmed   bool    `json:"confirmed"`
	Deactivated bool    `json:"deactivated"`
	Notes       *string `json:"notes"`
	HasLabel    bool    `json:"has_label"`
}
