package dtos

type DashboardResponse struct {
	Assets     AssetStats     `json:"assets"`
	Inventory  InventoryStats `json:"inventory"`
	Categories []CategoryStat `json:"categories"`
	Cities     []CityStat     `json:"cities"`
}

type AssetStats struct {
	Total      int `json:"total"`
	Active     int `json:"active"`
	Inactive   int `json:"inactive"`
	WrittenOff int `json:"written_off"`
}

type InventoryStats struct {
	OpenPeriod   *OpenPeriodStat   `json:"open_period"`
	LastClosed   *ClosedPeriodStat `json:"last_closed"`
	TotalPeriods int               `json:"total_periods"`
}

type OpenPeriodStat struct {
	ID          string  `json:"id"`
	PeriodYear  int     `json:"period_year"`
	PeriodMonth int     `json:"period_month"`
	PeriodDay 	int 	`json:"period_day"`
	Reviewed    int     `json:"reviewed"`
	Total       int     `json:"total"`
	Percentage  float64 `json:"percentage"`
}

type ClosedPeriodStat struct {
	PeriodYear  int `json:"period_year"`
	PeriodMonth int `json:"period_month"`
	PeriodDay 	int `json:"period_day"`
}

type CategoryStat struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type CityStat struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}
