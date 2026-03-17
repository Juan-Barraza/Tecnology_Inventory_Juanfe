package repository

import (
	"database/sql"

	response "inventory-juanfe/dtos/response"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetAssetStats() (response.AssetStats, error) {
	var stats response.AssetStats
	err := r.db.QueryRow(`
		SELECT
			COUNT(*)                                                      AS total,
			COUNT(*) FILTER (WHERE logical_status = 'active')            AS active,
			COUNT(*) FILTER (WHERE logical_status = 'inactive')          AS inactive,
			COUNT(*) FILTER (WHERE logical_status = 'written_off')       AS written_off
		FROM assets`).Scan(
		&stats.Total, &stats.Active, &stats.Inactive, &stats.WrittenOff,
	)
	return stats, err
}

func (r *DashboardRepository) GetCategoryStats() ([]response.CategoryStat, error) {
	rows, err := r.db.Query(`
		SELECT ac.name, COUNT(a.id) AS total
		FROM asset_categories ac
		LEFT JOIN assets a ON a.category_id = ac.id
		GROUP BY ac.name
		ORDER BY total DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []response.CategoryStat{}
	for rows.Next() {
		var s response.CategoryStat
		if err := rows.Scan(&s.Name, &s.Total); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *DashboardRepository) GetCityStats() ([]response.CityStat, error) {
	rows, err := r.db.Query(`
		SELECT c.name, COUNT(a.id) AS total
		FROM cities c
		LEFT JOIN assets a ON a.city_id = c.id
		GROUP BY c.name
		ORDER BY total DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []response.CityStat{}
	for rows.Next() {
		var s response.CityStat
		if err := rows.Scan(&s.Name, &s.Total); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *DashboardRepository) GetInventoryStats() (response.InventoryStats, error) {
	var stats response.InventoryStats

	// Total períodos
	err := r.db.QueryRow(`SELECT COUNT(*) FROM inventory_periods`).Scan(&stats.TotalPeriods)
	if err != nil {
		return stats, err
	}

	// Período abierto
	var openID, openYear, openMonth sql.NullString
	err = r.db.QueryRow(`
		SELECT id, period_year::text, period_month::text
		FROM inventory_periods
		WHERE status = 'open'
		LIMIT 1`).Scan(&openID, &openYear, &openMonth)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	if openID.Valid {
		var reviewed, total int
		_ = r.db.QueryRow(`
			SELECT
				COUNT(DISTINCT ir.asset_id) AS reviewed,
				COUNT(DISTINCT a.id)        AS total
			FROM assets a
			LEFT JOIN inventory_records ir
				ON ir.asset_id = a.id AND ir.period_id = $1
			WHERE a.logical_status = 'active'
				OR (a.logical_status = 'written_off' AND ir.id IS NOT NULL)`,
			openID.String,
		).Scan(&reviewed, &total)

		var pct float64
		if total > 0 {
			pct = float64(reviewed) / float64(total) * 100
		}

		var yr, mo int
		_ = r.db.QueryRow(`
			SELECT period_year, period_month FROM inventory_periods WHERE id = $1`,
			openID.String,
		).Scan(&yr, &mo)

		stats.OpenPeriod = &response.OpenPeriodStat{
			ID:          openID.String,
			PeriodYear:  yr,
			PeriodMonth: mo,
			Reviewed:    reviewed,
			Total:       total,
			Percentage:  pct,
		}
	}

	// Último período cerrado
	var closedYear, closedMonth int
	err = r.db.QueryRow(`
		SELECT period_year, period_month
		FROM inventory_periods
		WHERE status = 'closed'
		ORDER BY period_year DESC, period_month DESC
		LIMIT 1`).Scan(&closedYear, &closedMonth)
	if err == nil {
		stats.LastClosed = &response.ClosedPeriodStat{
			PeriodYear:  closedYear,
			PeriodMonth: closedMonth,
		}
	}

	return stats, nil
}
