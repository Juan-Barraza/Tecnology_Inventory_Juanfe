package dtos

type StatusHistoryResponse struct {
	ID             string  `json:"id"`
	AssetID        string  `json:"asset_id"`
	PreviousStatus *string `json:"previous_status"`
	NewStatus      string  `json:"new_status"`
	Notes          *string `json:"notes"`
	RecordedBy     string  `json:"recorded_by"`
	RecordedByName string  `json:"recorded_by_name"`
	CreatedAt      string  `json:"created_at"`
}
