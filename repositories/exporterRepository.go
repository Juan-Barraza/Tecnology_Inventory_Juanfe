package repository

import (
	"database/sql"
	dtos "inventory-juanfe/dtos/response"
	"inventory-juanfe/models"
)

type ExporterRepository struct {
	db *sql.DB
}

func NewExporterRepository(db *sql.DB) *ExporterRepository {
	return &ExporterRepository{db: db}
}

func (r *ExporterRepository) GetAssetsWithDate(year, month, day int, ownerId string) ([]models.AssetExport, error) {
	assets := []models.AssetExport{}
	query := `
		 SELECT 
			a.code,
			a.description,
			a.historical_cost,
			a.activation_date,
			a.logical_status,
			a.physical_status,
			ac.name as category,
			COALESCE(ar.name, '') as area,
			c.name as city,
			COALESCE(r.responsible_name, ''),
			COALESCE(r.responsible_position, ''),
			p.period_year,
			p.period_month,
			p.period_day,
			asac.code as accounting_group,
			acg.account_code as sub_code,
			ir.confirmed ,
			ir.deactivated,
			ir.has_label  
		FROM assets a 
		JOIN asset_categories ac on ac.id = a.category_id
		LEFT JOIN areas ar on ar.id = a.area_id 
		JOIN cities c on c.id = a.city_id 
		LEFT JOIN assignments r on a.id = r.asset_id
		JOIN inventory_records ir on ir.asset_id = a.id 
		JOIN inventory_periods p on ir.period_id = p.id
		JOIN asset_accounts acg on acg.id = a.asset_account_id
		JOIN accounting_groups asac on acg.accounting_group_id  = asac.id
		WHERE a.owner_id = $4 AND 
		(p.period_year = $1 AND p.period_month = $2 AND p.period_day = $3)
		ORDER BY a.activation_date desc
	`
	rows, err := r.db.Query(query, year, month, day, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.AssetExport
		if err := rows.Scan(
			&a.Code,
			&a.Description,
			&a.HistoricalCost,
			&a.ActivationDate,
			&a.LogicalStatus,
			&a.PhysicalStatus,
			&a.CategoryName,
			&a.AreaName,
			&a.CityName,
			&a.ResponsibleName,
			&a.ResponsiblePosition,
			&a.PeriodYear,
			&a.PeriodMonth,
			&a.PeriodDay,
			&a.AccountCodeGroup,
			&a.SubCode,
			&a.Confirmed,
			&a.Deactivated,
			&a.HasLabel,
		); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *ExporterRepository) GetAssetsToExport(ownerId string) ([]models.AssetExport, error) {
	assests := []models.AssetExport{}
	query := `
			SELECT 
				a.code,
				a.description,
				a.historical_cost,
				a.activation_date,
				a.logical_status,
				a.physical_status,
				ac.name as category,
				COALESCE(ar.name, '') as area,
				c.name as city,
				COALESCE(r.responsible_name, ''),
				COALESCE(r.responsible_position,''),
				asac.code as accounting_group,
				acg.account_code as sub_code 
			FROM assets a 
			JOIN asset_categories ac on ac.id = a.category_id
			LEFT JOIN areas ar on ar.id = a.area_id 
			JOIN cities c on c.id = a.city_id 
			LEFT JOIN assignments r on a.id = r.asset_id
			JOIN asset_accounts acg on acg.id = a.asset_account_id
			JOIN accounting_groups asac on acg.accounting_group_id  = asac.id
			WHERE a.owner_id = $1
			ORDER BY a.activation_date desc 
	`
	rows, err := r.db.Query(query, ownerId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var a models.AssetExport
		if err := rows.Scan(
			&a.Code,
			&a.Description,
			&a.HistoricalCost,
			&a.ActivationDate,
			&a.LogicalStatus,
			&a.PhysicalStatus,
			&a.CategoryName,
			&a.AreaName,
			&a.CityName,
			&a.ResponsibleName,
			&a.ResponsiblePosition,
			&a.AccountCodeGroup,
			&a.SubCode,
		); err != nil {
			return nil, err
		}

		assests = append(assests, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assests, nil
}

// posible error de direccion de memoria
func (r *ExporterRepository) CountAssetsConfirmatedAndDesactivated(year, month, day int, ownerId string) (*dtos.CounterAssetsToExport, error) {
	var responseQ dtos.CounterAssetsToExport
	query := `
		SELECT 
			COALESCE(SUM(ir.confirmed::int), 0) as total_confirmed,
			COALESCE(SUM(ir.deactivated::int),0) as total_desactivated,
			COALESCE(SUM(ir.has_label::int),0) as total_has_label,
			COALESCE(SUM(case when ir.has_label = false then 1 else 0 end), 0) AS total_without_label
		FROM inventory_records ir
		JOIN inventory_periods p ON ir.period_id = p.id
		JOIN assets a ON a.id = ir.asset_id
		WHERE a.owner_id = $4
			AND (p.period_year = $1 AND p.period_month = $2 AND p.period_day = $3)
	`
	err := r.db.QueryRow(query, year, month, day, ownerId).
		Scan(
			&responseQ.TotalConfirmated,
			&responseQ.TotalDesactivated,
			&responseQ.TotalWithLabel,
			&responseQ.TotalWithoutLabel,
		)
	if err != nil {
		return nil, err
	}

	return &responseQ, nil
}
