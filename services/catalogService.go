package services

import (
	repository "inventory-juanfe/repositories"
)

type CatalogService struct {
	cityRepo      *repository.CityRepository
	areaRepo      *repository.AreaRepository
	categoryRepo  *repository.CategoryRepository
	acctGroupRepo *repository.AccountingGroupRepository
}

func NewCatalogService(
	cityRepo *repository.CityRepository,
	areaRepo *repository.AreaRepository,
	categoryRepo *repository.CategoryRepository,
	acctGroupRepo *repository.AccountingGroupRepository,
) *CatalogService {
	return &CatalogService{
		cityRepo:      cityRepo,
		areaRepo:      areaRepo,
		categoryRepo:  categoryRepo,
		acctGroupRepo: acctGroupRepo,
	}
}

func (s *CatalogService) ListCities() (interface{}, error) {
	return s.cityRepo.FindAll()
}

func (s *CatalogService) ListAreas() (interface{}, error) {
	return s.areaRepo.FindAll()
}

func (s *CatalogService) ListCategories() (interface{}, error) {
	return s.categoryRepo.FindAll()
}

func (s *CatalogService) ListAccountingGroups() (interface{}, error) {
	return s.acctGroupRepo.FindAllWithAccounts()
}

func (s *CatalogService) UpdateAccountingGroup(id int, name string) error {
	return s.acctGroupRepo.UpdateName(id, name)
}
