package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	dtos "inventory-juanfe/dtos/request"
	response "inventory-juanfe/dtos/response"
	"inventory-juanfe/models"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/utils"
)

type AssetService struct {
	assetRepo   *repository.AssetRepository
	historyRepo *repository.StatusHistoryRepository
	assignRepo  *repository.AssignmentRepository
}

func NewAssetService(
	assetRepo *repository.AssetRepository,
	historyRepo *repository.StatusHistoryRepository,
	assignRepo *repository.AssignmentRepository,
) *AssetService {
	return &AssetService{
		assetRepo:   assetRepo,
		historyRepo: historyRepo,
		assignRepo:  assignRepo,
	}
}

func (s *AssetService) List(f dtos.AssetFilter, userId string) ([]dtos.AssetResponse, int, error) {
	assets, total, err := s.assetRepo.FindAll(f, userId)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]dtos.AssetResponse, len(assets))
	for i, a := range assets {
		resp[i] = toAssetResponse(a)
	}
	return resp, total, nil
}

func (s *AssetService) GetByID(id, userId string) (*dtos.AssetResponse, error) {
	a, err := s.assetRepo.FindByID(id, userId)
	if err != nil {
		return nil, fmt.Errorf("activo no encontrado")
	}
	if a == nil {
		return nil, nil
	}
	r := toAssetResponse(*a)
	return &r, nil
}

// Create builds and persists a new asset.
// Field-level validations are handled by the handler via utils.ValidateCreateAsset.
func (s *AssetService) Create(req dtos.CreateAssetRequest, userID string) (*dtos.AssetResponse, error) {
	existing, err := s.assetRepo.FindByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("asset code already exists")
	}

	actDate := time.Now()
	if req.ActivationDate != "" {
		if parsed, err := time.Parse(time.DateOnly, req.ActivationDate); err == nil {
			actDate = parsed
		}
	}

	physical := models.PhysicalStatus(req.PhysicalStatus)
	if !utils.IsValidPhysicalStatus(physical) {
		physical = models.PhysicalStatusOptimal
	}

	asset := &models.Asset{
		ID:             uuid.NewString(),
		Code:           req.Code,
		Description:    req.Description,
		OwnerId:        userID,
		Owner:          req.Owner,
		CategoryID:     req.CategoryID,
		AssetAccountID: req.AssetAccountID,
		CityID:         req.CityID,
		AreaID:         req.AreaID,
		HistoricalCost: req.HistoricalCost,
		ActivationDate: actDate,
		LogicalStatus:  models.LogicalStatusActive,
		PhysicalStatus: physical,
	}

	if err := s.assetRepo.Create(asset); err != nil {
		return nil, err
	}

	// primer registro en status_history
	_ = s.historyRepo.Create(&models.StatusHistory{
		ID:         uuid.NewString(),
		AssetID:    asset.ID,
		NewStatus:  models.LogicalStatusActive,
		RecordedBy: userID,
	})

	return s.GetByID(asset.ID, userID)
}

func (s *AssetService) Update(id string, req dtos.UpdateAssetRequest, userId string) (*dtos.AssetResponse, error) {
	detail, err := s.assetRepo.FindByID(id, userId)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, nil
	}

	a := &detail.Asset

	if req.Code != nil {
		a.Code = req.Code
	}

	if req.Description != nil {
		a.Description = *req.Description
	}
	if req.CategoryID != nil {
		a.CategoryID = *req.CategoryID
	}
	if req.AssetAccountID != nil {
		a.AssetAccountID = *req.AssetAccountID
	}
	if req.CityID != nil {
		a.CityID = *req.CityID
	}
	if req.AreaID != nil {
		a.AreaID = req.AreaID
	}
	if req.HistoricalCost != nil {
		a.HistoricalCost = req.HistoricalCost
	}
	if req.PhysicalStatus != nil {
		a.PhysicalStatus = models.PhysicalStatus(*req.PhysicalStatus)
	}
	if req.Owner != nil {
		a.Owner = req.Owner
	}

	if err := s.assetRepo.Update(a); err != nil {
		return nil, fmt.Errorf("error al acualizar")
	}
	return s.GetByID(id, userId)
}

// ChangeStatus updates logical/physical status.
// Field-level validations are handled by the handler via utils.ValidateUpdateAssetStatus.
func (s *AssetService) ChangeStatus(id string, req dtos.UpdateAssetStatusRequest, userID string) (*dtos.AssetResponse, error) {
	detail, err := s.assetRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, nil
	}

	prev := detail.LogicalStatus
	newLogical := prev
	newPhysical := detail.PhysicalStatus

	if req.LogicalStatus != nil {
		newLogical = models.LogicalStatus(*req.LogicalStatus)
	}
	if req.PhysicalStatus != nil {
		newPhysical = models.PhysicalStatus(*req.PhysicalStatus)
	}

	if err := s.assetRepo.UpdateStatus(id, newLogical, newPhysical); err != nil {
		return nil, err
	}

	// solo escribe historial si cambió el estado lógico
	if newLogical != prev {
		_ = s.historyRepo.Create(&models.StatusHistory{
			ID:             uuid.NewString(),
			AssetID:        id,
			PreviousStatus: &prev,
			NewStatus:      newLogical,
			Notes:          req.Notes,
			RecordedBy:     userID,
		})

		// si se da de baja, cerrar el assignment activo
		if newLogical == models.LogicalStatusWrittenOff {
			now := time.Now()
			reason := "asset written off"
			_ = s.assignRepo.WriteOff(id, now, &reason)
		}
	}

	return s.GetByID(id, userID)
}

func (s *AssetService) GetHistory(assetID string) ([]response.StatusHistoryResponse, error) {
	records, err := s.historyRepo.FindByAssetID(assetID)
	if err != nil {
		return nil, err
	}

	result := make([]response.StatusHistoryResponse, len(records))
	for i, h := range records {
		var prev *string
		if h.PreviousStatus != nil {
			s := string(*h.PreviousStatus)
			prev = &s
		}
		result[i] = response.StatusHistoryResponse{
			ID:             h.ID,
			AssetID:        h.AssetID,
			PreviousStatus: prev,
			NewStatus:      string(h.NewStatus),
			Notes:          h.Notes,
			RecordedBy:     h.RecordedBy,
			RecordedByName: h.RecordedByName,
			CreatedAt:      h.CreatedAt.Format(time.RFC3339),
		}
	}
	return result, nil
}

// ── helpers ───────────────────────────────────────────────────

func toAssetResponse(a models.AssetDetail) dtos.AssetResponse {
	return dtos.AssetResponse{
		ID:                  a.ID,
		Code:                a.Code,
		Description:         a.Description,
		Owner:               a.Owner,
		Category:            a.CategoryName,
		AccountingGroupName: a.AccountingGroupName,
		AccountingGroupCode: a.AccountingGroupCode,
		AccountCode:         a.AccountCode,
		OpenLedger:          a.OpenLedger,
		City:                a.CityName,
		CityId:              a.CityID,
		Area:                a.AreaName,
		HistoricalCost:      a.HistoricalCost,
		ActivationDate:      a.ActivationDate.Format(time.DateOnly),
		LogicalStatus:       string(a.LogicalStatus),
		PhysicalStatus:      string(a.PhysicalStatus),
		CreatedAt:           a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           a.UpdatedAt.Format(time.RFC3339),
	}
}
