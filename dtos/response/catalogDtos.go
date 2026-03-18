package dtos

type CityResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
}

type AreaResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type AssetCategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AssetAccountItemResponse struct {
	ID          int     `json:"id"`
	AccountCode int64   `json:"account_code"`
	OpenLedger  *string `json:"open_ledger"`
}

type AccountingGroupResponse struct {
	ID       int                        `json:"id"`
	Code     int64                      `json:"code"`
	Name     string                     `json:"name"`
	Accounts []AssetAccountItemResponse `json:"accounts"`
}
