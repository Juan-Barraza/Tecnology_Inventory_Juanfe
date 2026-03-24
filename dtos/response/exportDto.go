package dtos

type CounterAssetsToExport struct {
	TotalConfirmated  int64 `json:"total_confirmated"`
	TotalDesactivated int64 `json:"total_desactivated"`
	TotalWithLabel    int64 `json:"total_with_label"`
	TotalWithoutLabel int64 `json:"total_without_label"`
}
