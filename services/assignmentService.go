package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	dtos "inventory-juanfe/dtos/request"
	"inventory-juanfe/models"
	repository "inventory-juanfe/repositories"
)

type AssignmentService struct {
	assignRepo *repository.AssignmentRepository
	assetRepo  *repository.AssetRepository
}

func NewAssignmentService(
	assignRepo *repository.AssignmentRepository,
	assetRepo *repository.AssetRepository,
) *AssignmentService {
	return &AssignmentService{
		assignRepo: assignRepo,
		assetRepo:  assetRepo,
	}
}

func (s *AssignmentService) GetByAsset(assetID string) ([]dtos.AssignmentResponse, error) {
	assignments, err := s.assignRepo.FindByAssetID(assetID)
	if err != nil {
		return nil, err
	}

	resp := make([]dtos.AssignmentResponse, len(assignments))
	for i, a := range assignments {
		resp[i] = toAssignmentResponse(a)
	}
	return resp, nil
}

// Create persists a new assignment.
// Field-level validations are handled by the handler via utils.ValidateCreateAssignment.
func (s *AssignmentService) Create(req dtos.CreateAssignmentRequest, userID string) (*dtos.AssignmentResponse, error) {
	asset, err := s.assetRepo.FindByID(req.AssetID, userID)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.New("asset not found")
	}
	if asset.LogicalStatus != models.LogicalStatusActive {
		return nil, errors.New("only active assets can be assigned")
	}

	current, err := s.assignRepo.FindActiveByAssetID(req.AssetID)
	if err != nil {
		return nil, err
	}
	if current != nil {
		return nil, errors.New("asset already has an active assignment — release it first")
	}

	assignedAt, _ := time.Parse(time.DateOnly, req.AssignedAt) // ya validado en handler

	a := &models.Assignment{
		ID:                  uuid.NewString(),
		AssetID:             req.AssetID,
		ResponsibleName:     req.ResponsibleName,
		ResponsiblePosition: req.ResponsiblePosition,
		AssignedAt:          assignedAt,
		Status:              models.AssignmentStatusActive,
		CreatedBy:           userID,
	}

	if err := s.assignRepo.Create(a); err != nil {
		return nil, err
	}

	detail, err := s.assignRepo.FindActiveByAssetID(req.AssetID)
	if err != nil || detail == nil {
		return nil, errors.New("assignment created but could not retrieve it")
	}

	r := toAssignmentResponse(*detail)
	return &r, nil
}

// Release deactivates an assignment.
// Field-level validations are handled by the handler via utils.ValidateReleaseAssignment.
func (s *AssignmentService) Release(id string, req dtos.ReleaseAssignmentRequest) error {
	deactivatedAt, _ := time.Parse(time.DateOnly, req.DeactivatedAt) // ya validado en handler
	return s.assignRepo.Release(id, deactivatedAt, req.DeactivationReason)
}

// ── helper ────────────────────────────────────────────────────

func toAssignmentResponse(a models.AssignmentDetail) dtos.AssignmentResponse {
	r := dtos.AssignmentResponse{
		ID:                  a.ID,
		AssetID:             a.AssetID,
		AssetCode:           a.AssetCode,
		AssetDescription:    a.AssetDescription,
		ResponsibleName:     a.ResponsibleName,
		ResponsiblePosition: a.ResponsiblePosition,
		AssignedAt:          a.AssignedAt.Format(time.DateOnly),
		DeactivationReason:  a.DeactivationReason,
		Status:              string(a.Status),
		CreatedByName:       a.CreatedByName,
		CreatedAt:           a.CreatedAt.Format(time.RFC3339),
	}
	if a.DeactivatedAt != nil {
		s := a.DeactivatedAt.Format(time.DateOnly)
		r.DeactivatedAt = &s
	}
	return r
}
