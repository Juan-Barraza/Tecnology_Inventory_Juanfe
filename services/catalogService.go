package services

import (
	response "inventory-juanfe/dtos/response"
	"inventory-juanfe/models"
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

func (s *CatalogService) ListCities() ([]response.CityResponse, error) {
	cities, err := s.cityRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var res []response.CityResponse
	for _, c := range cities {
		res = append(res, toCityResponse(c))
	}
	return res, nil
}

func (s *CatalogService) ListAreas() ([]response.AreaResponse, error) {
	areas, err := s.areaRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var res []response.AreaResponse
	for _, a := range areas {
		res = append(res, toAreaResponse(a))
	}
	return res, nil
}

func (s *CatalogService) ListCategories() ([]response.AssetCategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var res []response.AssetCategoryResponse
	for _, c := range categories {
		res = append(res, toCategoryResponse(c))
	}
	return res, nil
}

func (s *CatalogService) ListAccountingGroups() ([]response.AccountingGroupResponse, error) {
	groups, err := s.acctGroupRepo.FindAllWithAccounts()
	if err != nil {
		return nil, err
	}
	var res []response.AccountingGroupResponse
	for _, g := range groups {
		res = append(res, toAccountingGroupResponse(g))
	}
	return res, nil
}

func (s *CatalogService) UpdateAccountingGroup(id int, name string) error {
	return s.acctGroupRepo.UpdateName(id, name)
}

func toCityResponse(c models.City) response.CityResponse {
	return response.CityResponse{ID: c.ID, Name: c.Name, Department: c.Department}
}
func toAreaResponse(a models.Area) response.AreaResponse {
	return response.AreaResponse{ID: a.ID, Name: a.Name, Description: a.Description}
}
func toCategoryResponse(c models.AssetCategory) response.AssetCategoryResponse {
	return response.AssetCategoryResponse{ID: c.ID, Name: c.Name}
}
func toAccountingGroupResponse(g models.AccountingGroupWithAccounts) response.AccountingGroupResponse {
	accounts := make([]response.AssetAccountItemResponse, len(g.Accounts))
	for i, a := range g.Accounts {
		accounts[i] = response.AssetAccountItemResponse{
			ID:          a.ID,
			AccountCode: a.AccountCode,
			OpenLedger:  a.OpenLedger,
		}
	}
	return response.AccountingGroupResponse{
		ID:       g.ID,
		Code:     g.Code,
		Name:     g.Name,
		Accounts: accounts,
	}
}
