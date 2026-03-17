package dtos

type CreateAssignmentRequest struct {
	AssetID             string  `json:"asset_id"`
	ResponsibleName     *string `json:"responsible_name"`
	ResponsiblePosition *string `json:"responsible_position"`
	AssignedAt          string  `json:"assigned_at"` // "YYYY-MM-DD"
}

type ReleaseAssignmentRequest struct {
	DeactivatedAt      string  `json:"deactivated_at"` // "YYYY-MM-DD"
	DeactivationReason *string `json:"deactivation_reason"`
}

type AssignmentResponse struct {
	ID                  string  `json:"id"`
	AssetID             string  `json:"asset_id"`
	AssetCode           string  `json:"asset_code"`
	AssetDescription    string  `json:"asset_description"`
	ResponsibleName     *string `json:"responsible_name"`
	ResponsiblePosition *string `json:"responsible_position"`
	AssignedAt          string  `json:"assigned_at"`
	DeactivatedAt       *string `json:"deactivated_at"`
	DeactivationReason  *string `json:"deactivation_reason"`
	Status              string  `json:"status"`
	CreatedByName       string  `json:"created_by_name"`
	CreatedAt           string  `json:"created_at"`
}
