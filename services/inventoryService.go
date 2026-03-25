package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	dtos "inventory-juanfe/dtos/request"
	response "inventory-juanfe/dtos/response"
	"inventory-juanfe/models"
	repository "inventory-juanfe/repositories"
)

type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
	assetRepo     *repository.AssetRepository
	historyRepo   *repository.StatusHistoryRepository
}

func NewInventoryService(
	inventoryRepo *repository.InventoryRepository,
	assetRepo *repository.AssetRepository,
	historyRepo *repository.StatusHistoryRepository,
) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		assetRepo:     assetRepo,
		historyRepo:   historyRepo,
	}
}

func (s *InventoryService) ListPeriods(userId string) ([]response.InventoryPeriodResponse, error) {
	periods, err := s.inventoryRepo.FindAllPeriods(userId)
	if err != nil {
		return nil, err
	}
	result := make([]response.InventoryPeriodResponse, len(periods))
	for i, p := range periods {
		result[i] = toPeriodResponse(p)
	}
	return result, nil
}

func (s *InventoryService) GetPeriod(id string, userID string) (*response.InventoryPeriodResponse, error) {
	p, err := s.inventoryRepo.FindPeriodByID(id, userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	r := toPeriodResponse(*p)
	return &r, nil
}

func (s *InventoryService) CreatePeriod(year, month, day int, userID string) (*response.InventoryPeriodResponse, error) {
	open, err := s.inventoryRepo.FindOpenPeriod(userID)
	if err != nil {
		return nil, err
	}
	if open != nil {
		return nil, errors.New("there is already an open period — close it before creating a new one")
	}

	p := &models.InventoryPeriod{
		ID:          uuid.NewString(),
		PeriodYear:  year,
		PeriodMonth: month,
		PeriodDay:   day,
		CreatedBy:   userID,
	}

	if err := s.inventoryRepo.CreatePeriod(p); err != nil {
		return nil, err
	}

	created, err := s.inventoryRepo.FindPeriodByID(p.ID, userID)
	if err != nil {
		return nil, err
	}
	r := toPeriodResponse(*created)
	return &r, nil
}

func (s *InventoryService) ClosePeriod(id, userID string) error {
	period, err := s.inventoryRepo.FindPeriodByID(id, userID)
	if err != nil {
		return err
	}
	if period == nil {
		return errors.New("period not found")
	}
	if period.Status == models.PeriodStatusClosed {
		return errors.New("period is already closed")
	}
	return s.inventoryRepo.ClosePeriod(id, time.Now(), userID)
}

func (s *InventoryService) GetRecords(periodID string, userId string) ([]response.InventoryRecordDetailResponse, error) {
	records, err := s.inventoryRepo.FindRecordsByPeriod(periodID, userId)
	if err != nil {
		return nil, err
	}
	result := make([]response.InventoryRecordDetailResponse, len(records))
	for i, r := range records {
		result[i] = response.InventoryRecordDetailResponse{
			ID:               r.ID,
			PeriodID:         r.PeriodID,
			AssetID:          r.AssetID,
			AssetCode:        r.AssetCode,
			AssetDescription: r.AssetDescription,
			CategoryName:     r.CategoryName,
			CityName:         r.CityName,
			AreaName:         r.AreaName,
			LogicalStatus:    r.LogicalStatus,
			Confirmed:        r.Confirmed,
			Deactivated:      r.Deactivated,
			Notes:            r.Notes,
			RecordedByName:   r.RecordedByName,
			RecordedAt:       r.RecordedAt,
		}
	}
	return result, nil
}

func (s *InventoryService) RecordAsset(req dtos.RecordAssetRequest, userID string) error {
	period, err := s.inventoryRepo.FindPeriodByID(req.PeriodID, userID)
	if err != nil {
		return err
	}
	if period == nil {
		return errors.New("period not found")
	}
	if period.Status == models.PeriodStatusClosed {
		return errors.New("cannot modify a closed period")
	}

	asset, err := s.assetRepo.FindByID(req.AssetID, userID)
	if err != nil {
		return err
	}
	if asset == nil {
		return errors.New("asset not found")
	}

	rec := &models.InventoryRecord{
		ID:          uuid.NewString(),
		PeriodID:    req.PeriodID,
		AssetID:     req.AssetID,
		Confirmed:   req.Confirmed,
		Deactivated: req.Deactivated,
		Notes:       req.Notes,
		HasLabel:    req.HasLabel,
		RecordedBy:  userID,
	}

	if err := s.inventoryRepo.UpsertRecord(rec); err != nil {
		return err
	}

	if req.Deactivated {
		prev := asset.LogicalStatus
		newStatus := models.LogicalStatusWrittenOff

		if err := s.assetRepo.UpdateStatus(req.AssetID, newStatus, asset.PhysicalStatus); err != nil {
			return err
		}

		_ = s.historyRepo.Create(&models.StatusHistory{
			ID:             uuid.NewString(),
			AssetID:        req.AssetID,
			PreviousStatus: &prev,
			NewStatus:      newStatus,
			Notes:          req.Notes,
			RecordedBy:     userID,
		})
	}

	return nil
}

func (s *InventoryService) GetPeriodAssets(periodID string, userId string) ([]response.AssetInventoryStatusResponse, error) {
	period, err := s.inventoryRepo.FindPeriodByID(periodID, userId)
	if err != nil {
		return nil, err
	}
	if period == nil {
		return nil, errors.New("period not found")
	}

	assets, err := s.inventoryRepo.FindAssetsWithPeriodStatus(periodID, userId)
	if err != nil {
		return nil, err
	}

	result := make([]response.AssetInventoryStatusResponse, len(assets))
	for i, a := range assets {
		result[i] = response.AssetInventoryStatusResponse{
			AssetID:          a.AssetID,
			AssetCode:        a.AssetCode,
			AssetDescription: a.AssetDescription,
			CategoryName:     a.CategoryName,
			CityName:         a.CityName,
			AreaName:         a.AreaName,
			RecordID:         a.RecordID,
			Confirmed:        a.Confirmed,
			Deactivated:      a.Deactivated,
			Notes:            a.Notes,
			RecordedAt:       a.RecordedAt,
			ActivationDate:   a.ActivationDate.Format("2006-01-02"),
			HasLabel:         a.HasLabel,
		}
	}
	return result, nil
}

func (s *InventoryService) GetProgress(periodID, userId string) (*response.PeriodProgressResponse, error) {
	total, reviewed, err := s.inventoryRepo.CountRecords(periodID, userId)
	if err != nil {
		return nil, err
	}

	pending := total - reviewed
	var pct float64
	if total > 0 {
		pct = float64(reviewed) / float64(total) * 100
	}

	return &response.PeriodProgressResponse{
		Total:      total,
		Reviewed:   reviewed,
		Pending:    pending,
		Percentage: pct,
	}, nil
}

// ── helpers ───────────────────────────────────────────────────

func toPeriodResponse(p models.InventoryPeriod) response.InventoryPeriodResponse {
	return response.InventoryPeriodResponse{
		ID:          p.ID,
		PeriodYear:  p.PeriodYear,
		PeriodMonth: p.PeriodMonth,
		PeriodDay:   p.PeriodDay,
		Status:      string(p.Status),
		CreatedBy:   p.CreatedBy,
		ClosedAt:    p.ClosedAt,
		CreatedAt:   p.CreatedAt,
	}
}
