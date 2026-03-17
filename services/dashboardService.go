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

func (s *DashboardService) GetDashboard() (*response.DashboardResponse, error) {
	assets, err := s.repo.GetAssetStats()
	if err != nil {
		return nil, err
	}

	inventory, err := s.repo.GetInventoryStats()
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.GetCategoryStats()
	if err != nil {
		return nil, err
	}

	cities, err := s.repo.GetCityStats()
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
