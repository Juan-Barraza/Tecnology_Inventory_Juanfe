package services

import (
	response "inventory-juanfe/dtos/response"
	repository "inventory-juanfe/repositories"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetDashboard(userId string) (*response.DashboardResponse, error) {
	assets, err := s.repo.GetAssetStats(userId)
	if err != nil {
		return nil, err
	}

	inventory, err := s.repo.GetInventoryStats(userId)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.GetCategoryStats(userId)
	if err != nil {
		return nil, err
	}

	cities, err := s.repo.GetCityStats(userId)
	if err != nil {
		return nil, err
	}

	return &response.DashboardResponse{
		Assets:     assets,
		Inventory:  inventory,
		Categories: categories,
		Cities:     cities,
	}, nil
}
