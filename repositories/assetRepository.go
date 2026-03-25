package repository

import (
	"database/sql"
	"fmt"
	"strings"

	dtos "inventory-juanfe/dtos/request"
	"inventory-juanfe/models"
)

type AssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

const assetSelectBase = `
    SELECT
        a.id, a.code, a.description,
        a.category_id, a.asset_account_id, a.city_id, a.area_id,
        a.historical_cost, a.activation_date,
        a.logical_status, a.physical_status,
        a.created_at, a.updated_at,
        ac.name   AS category_name,
        ag.name   AS accounting_group_name,
        ag.code   AS accounting_group_code,
        aa.account_code,
        aa.open_ledger,
        c.name    AS city_name,
        ar.name   AS area_name
    FROM assets a
    JOIN asset_categories  ac ON ac.id  = a.category_id
    JOIN asset_accounts    aa ON aa.id  = a.asset_account_id
    JOIN accounting_groups ag ON ag.id  = aa.accounting_group_id
    JOIN cities            c  ON c.id   = a.city_id
    LEFT JOIN areas        ar ON ar.id  = a.area_id`

func (r *AssetRepository) FindAll(f dtos.AssetFilter, ownerId string) ([]models.AssetDetail, int, error) {
	where, args := buildAssetWhere(f, ownerId)

	var total int
	if err := r.db.QueryRow(
		"SELECT COUNT(*) FROM assets a"+where, args...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit
	n := len(args) + 1

	query := fmt.Sprintf(
		"%s%s ORDER BY a.activation_date DESC LIMIT $%d OFFSET $%d",
		assetSelectBase, where, n, n+1,
	)
	args = append(args, f.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assets []models.AssetDetail
	for rows.Next() {
		var a models.AssetDetail
		if err := scanAssetDetail(rows, &a); err != nil {
			return nil, 0, err
		}
		assets = append(assets, a)
	}
	return assets, total, rows.Err()
}

func (r *AssetRepository) FindByID(id, owner_id string) (*models.AssetDetail, error) {
	row := r.db.QueryRow(assetSelectBase+" WHERE a.id = $1 AND a.owner_id = $2", id, owner_id)
	var a models.AssetDetail
	if err := scanAssetDetail(row, &a); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AssetRepository) FindByCode(code string) (*models.Asset, error) {
	var a models.Asset
	err := r.db.QueryRow(`
        SELECT id, code, description, category_id, asset_account_id,
            city_id, area_id, historical_cost, activation_date,
            logical_status, physical_status, created_at, updated_at
        FROM assets WHERE code = $1`, code,
	).Scan(
		&a.ID, &a.Code, &a.Description, &a.CategoryID, &a.AssetAccountID,
		&a.CityID, &a.AreaID, &a.HistoricalCost, &a.ActivationDate,
		&a.LogicalStatus, &a.PhysicalStatus, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *AssetRepository) Create(a *models.Asset) error {
	_, err := r.db.Exec(`
        INSERT INTO assets
            (id, code, description, owner_id, category_id, asset_account_id,
             city_id, area_id, historical_cost, activation_date,
             logical_status, physical_status)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		a.ID, a.Code, a.Description, a.OwnerId, a.CategoryID, a.AssetAccountID,
		a.CityID, a.AreaID, a.HistoricalCost, a.ActivationDate,
		a.LogicalStatus, a.PhysicalStatus,
	)
	return err
}

func (r *AssetRepository) Update(a *models.Asset) error {
	_, err := r.db.Exec(`
        UPDATE assets SET
			code 			 = $1,
            description      = $2,
            category_id      = $3,
            asset_account_id = $4,
            city_id          = $5,
            area_id          = $6,
            historical_cost  = $7,
            physical_status  = $8,
            logical_status   = $9
        WHERE id = $10`,
		a.Code, a.Description, a.CategoryID, a.AssetAccountID,
		a.CityID, a.AreaID, a.HistoricalCost,
		a.PhysicalStatus, a.LogicalStatus, a.ID,
	)
	return err
}

func (r *AssetRepository) UpdateStatus(id string, logical models.LogicalStatus, physical models.PhysicalStatus) error {
	_, err := r.db.Exec(`
        UPDATE assets SET logical_status = $1, physical_status = $2 WHERE id = $3`,
		logical, physical, id,
	)
	return err
}

func buildAssetWhere(f dtos.AssetFilter, ownerId string) (string, []interface{}) {
	var conds []string
	var args []interface{}
	n := 1

	if ownerId != "" {
		conds = append(conds, fmt.Sprintf("a.owner_id = $%d", n))
		args = append(args, ownerId)
		n++
	}
	if f.CityID != nil {
		conds = append(conds, fmt.Sprintf("a.city_id = $%d", n))
		args = append(args, *f.CityID)
		n++
	}
	if f.AreaID != nil {
		conds = append(conds, fmt.Sprintf("a.area_id = $%d", n))
		args = append(args, *f.AreaID)
		n++
	}
	if f.CategoryID != nil {
		conds = append(conds, fmt.Sprintf("a.category_id = $%d", n))
		args = append(args, *f.CategoryID)
		n++
	}
	if f.AssetAccountID != nil {
		conds = append(conds, fmt.Sprintf("a.asset_account_id = $%d", n))
		args = append(args, *f.AssetAccountID)
		n++
	}
	if f.LogicalStatus != nil {
		conds = append(conds, fmt.Sprintf("a.logical_status = $%d", n))
		args = append(args, *f.LogicalStatus)
		n++
	}
	if f.PhysicalStatus != nil {
		conds = append(conds, fmt.Sprintf("a.physical_status = $%d", n))
		args = append(args, *f.PhysicalStatus)
		n++
	}
	if f.From != nil && *f.From != "" {
		conds = append(conds, fmt.Sprintf("a.activation_date >= $%d", n))
		args = append(args, *f.From)
		n++
	}
	if f.To != nil && *f.To != "" {
		conds = append(conds, fmt.Sprintf("a.activation_date <= $%d", n))
		args = append(args, *f.To)
		n++
	}
	if f.Search != nil && *f.Search != "" {
		conds = append(conds, fmt.Sprintf(
			"(a.code ILIKE $%d OR a.description ILIKE $%d)", n, n,
		))
		args = append(args, "%"+*f.Search+"%")
		n++
	}

	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

type scanner interface {
	Scan(dest ...any) error
}

func scanAssetDetail(s scanner, a *models.AssetDetail) error {
	return s.Scan(
		&a.ID, &a.Code, &a.Description,
		&a.CategoryID, &a.AssetAccountID, &a.CityID, &a.AreaID,
		&a.HistoricalCost, &a.ActivationDate,
		&a.LogicalStatus, &a.PhysicalStatus,
		&a.CreatedAt, &a.UpdatedAt,
		&a.CategoryName,
		&a.AccountingGroupName, &a.AccountingGroupCode,
		&a.AccountCode, &a.OpenLedger,
		&a.CityName,
		&a.AreaName,
	)
}
