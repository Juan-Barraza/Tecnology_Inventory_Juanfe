package utils

import (
	"errors"
	"time"

	dtos "inventory-juanfe/dtos/request"
	"inventory-juanfe/models"
)

// ── Asset validators ──────────────────────────────────────────

func ValidateCreateAsset(req dtos.CreateAssetRequest) error {
	if req.Code == "" {
		return errors.New("code is required")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	if req.ActivationDate == "" {
		return errors.New("activation_date is required")
	}
	if _, err := time.Parse(time.DateOnly, req.ActivationDate); err != nil {
		return errors.New("activation_date must be YYYY-MM-DD")
	}
	return nil
}

func ValidateUpdateAssetStatus(req dtos.UpdateAssetStatusRequest) error {
	if req.LogicalStatus == nil && req.PhysicalStatus == nil {
		return errors.New("at least one of logical_status or physical_status is required")
	}
	if req.LogicalStatus != nil && !IsValidLogicalStatus(models.LogicalStatus(*req.LogicalStatus)) {
		return errors.New("invalid logical_status value")
	}
	if req.PhysicalStatus != nil && !IsValidPhysicalStatus(models.PhysicalStatus(*req.PhysicalStatus)) {
		return errors.New("invalid physical_status value")
	}
	return nil
}

// ── Assignment validators ─────────────────────────────────────

func ValidateCreateAssignment(req dtos.CreateAssignmentRequest) error {
	if req.AssetID == "" {
		return errors.New("asset_id is required")
	}
	if req.AssignedAt == "" {
		return errors.New("assigned_at is required")
	}
	if _, err := time.Parse(time.DateOnly, req.AssignedAt); err != nil {
		return errors.New("assigned_at must be YYYY-MM-DD")
	}
	return nil
}

func ValidateReleaseAssignment(req dtos.ReleaseAssignmentRequest) error {
	if req.DeactivatedAt == "" {
		return errors.New("deactivated_at is required")
	}
	if _, err := time.Parse(time.DateOnly, req.DeactivatedAt); err != nil {
		return errors.New("deactivated_at must be YYYY-MM-DD")
	}
	return nil
}

// ── Inventory validators ──────────────────────────────────────

func ValidateCreatePeriod(req dtos.CreatePeriodRequest) error {
	if req.PeriodMonth < 1 || req.PeriodMonth > 12 {
		return errors.New("period_month must be between 1 and 12")
	}
	if req.PeriodYear < 2000 {
		return errors.New("period_year must be >= 2000")
	}
	return nil
}

func ValidateRecordAsset(req dtos.RecordAssetRequest) error {
	if req.PeriodID == "" {
		return errors.New("period_id is required")
	}
	if req.AssetID == "" {
		return errors.New("asset_id is required")
	}
	return nil
}

// ── Catalog validators ────────────────────────────────────────

func ValidateUpdateAccountingGroup(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	return nil
}

// ── Status helpers (exportados) ───────────────────────────────

func IsValidLogicalStatus(s models.LogicalStatus) bool {
	switch s {
	case models.LogicalStatusActive,
		models.LogicalStatusInactive,
		models.LogicalStatusWrittenOff:
		return true
	}
	return false
}

func IsValidPhysicalStatus(s models.PhysicalStatus) bool {
	switch s {
	case models.PhysicalStatusOptimal,
		models.PhysicalStatusGood,
		models.PhysicalStatusFair,
		models.PhysicalStatusDeteriorated,
		models.PhysicalStatusOutOfService:
		return true
	}
	return false
}
