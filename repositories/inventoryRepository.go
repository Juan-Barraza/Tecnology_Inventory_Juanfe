package repository

import (
	"database/sql"
	"time"

	"inventory-juanfe/models"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) FindAllPeriods(created_by string) ([]models.InventoryPeriod, error) {
	rows, err := r.db.Query(`
        SELECT 
			id, 
			period_year, 
			period_month, 
			period_day, 
			status, 
			created_by, 
			closed_at, 
			created_at
        FROM inventory_periods
		WHERE created_by = $1
        ORDER BY period_year DESC, period_month DESC, period_day DESC`, created_by)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	periods := []models.InventoryPeriod{}
	for rows.Next() {
		var p models.InventoryPeriod
		if err := rows.Scan(
			&p.ID, &p.PeriodYear, &p.PeriodMonth,
			&p.PeriodDay,
			&p.Status, &p.CreatedBy, &p.ClosedAt, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		periods = append(periods, p)
	}
	return periods, rows.Err()
}

func (r *InventoryRepository) FindPeriodByID(id string, created_by string) (*models.InventoryPeriod, error) {
	var p models.InventoryPeriod
	err := r.db.QueryRow(`
        SELECT 
			id,
			period_year, 
			period_month, 
			period_day, 
			status, 
			created_by, 
			closed_at,
			created_at
        FROM inventory_periods
        WHERE id = $1 AND created_by = $2`, id, created_by).Scan(
		&p.ID, &p.PeriodYear, &p.PeriodMonth,
		&p.PeriodDay,
		&p.Status, &p.CreatedBy, &p.ClosedAt, &p.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *InventoryRepository) FindOpenPeriod(created_by string) (*models.InventoryPeriod, error) {
	var p models.InventoryPeriod
	err := r.db.QueryRow(`
        SELECT id, period_year, period_month, period_day, status, created_by, closed_at, created_at
        FROM inventory_periods
        WHERE status = 'open' AND created_by = $1
        LIMIT 1`, created_by).Scan(
		&p.ID, &p.PeriodYear, &p.PeriodMonth,
		&p.PeriodDay,
		&p.Status, &p.CreatedBy, &p.ClosedAt, &p.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *InventoryRepository) CreatePeriod(p *models.InventoryPeriod) error {
	_, err := r.db.Exec(`
        INSERT INTO inventory_periods (id, period_year, period_month, period_day, status, created_by)
        VALUES ($1, $2, $3, $4, 'open', $5)`,
		p.ID, p.PeriodYear, p.PeriodMonth, p.PeriodDay, p.CreatedBy,
	)
	return err
}

func (r *InventoryRepository) ClosePeriod(id string, closedAt time.Time, created_by string) error {
	_, err := r.db.Exec(`
        UPDATE inventory_periods
        SET status = 'closed', closed_at = $1
        WHERE id = $2 AND created_by = $3`,
		closedAt, id, created_by,
	)
	return err
}

func (r *InventoryRepository) FindRecordsByPeriod(periodID string, ownerId string) ([]models.InventoryRecordDetail, error) {
	rows, err := r.db.Query(`
        SELECT
            ir.id, ir.period_id, ir.asset_id,
            ir.confirmed, ir.deactivated, ir.notes,
            ir.recorded_by, ir.recorded_at,
            a.code        AS asset_code,
            a.description AS asset_description,
            ac.name       AS category_name,
            c.name        AS city_name,
            ar.name       AS area_name,
            a.logical_status,
            u.name        AS recorded_by_name
        FROM inventory_records ir
        JOIN assets          a  ON a.id  = ir.asset_id
        JOIN asset_categories ac ON ac.id = a.category_id
        JOIN cities          c  ON c.id  = a.city_id
        LEFT JOIN areas      ar ON ar.id = a.area_id
        JOIN users           u  ON u.id  = ir.recorded_by
        WHERE ir.period_id = $1 AND a.owner_id = $2
        ORDER BY ir.recorded_at DESC`, periodID, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []models.InventoryRecordDetail{}
	for rows.Next() {
		var rec models.InventoryRecordDetail
		if err := rows.Scan(
			&rec.ID, &rec.PeriodID, &rec.AssetID,
			&rec.Confirmed, &rec.Deactivated, &rec.Notes,
			&rec.RecordedBy, &rec.RecordedAt,
			&rec.AssetCode, &rec.AssetDescription,
			&rec.CategoryName, &rec.CityName, &rec.AreaName,
			&rec.LogicalStatus, &rec.RecordedByName,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

func (r *InventoryRepository) FindRecordByAsset(periodID, assetID string) (*models.InventoryRecord, error) {
	var rec models.InventoryRecord
	err := r.db.QueryRow(`
        SELECT id, period_id, asset_id, confirmed, deactivated, notes, recorded_by, recorded_at
        FROM inventory_records
        WHERE period_id = $1 AND asset_id = $2`,
		periodID, assetID,
	).Scan(
		&rec.ID, &rec.PeriodID, &rec.AssetID,
		&rec.Confirmed, &rec.Deactivated, &rec.Notes,
		&rec.RecordedBy, &rec.RecordedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &rec, err
}

func (r *InventoryRepository) UpsertRecord(rec *models.InventoryRecord) error {
	_, err := r.db.Exec(`
        INSERT INTO inventory_records
            (id, period_id, asset_id, confirmed, deactivated, has_label, notes, recorded_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (period_id, asset_id) DO UPDATE SET
            confirmed   = EXCLUDED.confirmed,
            deactivated = EXCLUDED.deactivated,
			has_label 	= EXCLUDED.has_label,
            notes       = EXCLUDED.notes,
            recorded_by = EXCLUDED.recorded_by,
            recorded_at = NOW()`,
		rec.ID, rec.PeriodID, rec.AssetID,
		rec.Confirmed, rec.Deactivated, rec.HasLabel, rec.Notes, rec.RecordedBy,
	)
	return err
}

// CountRecords devuelve total de activos activos y cuántos fueron revisados en el período
func (r *InventoryRepository) CountRecords(periodID string, ownerId string) (total int, reviewed int, err error) {
	err = r.db.QueryRow(`
        SELECT
            COUNT(DISTINCT a.id)                                    AS total,
            COUNT(DISTINCT ir.asset_id)                             AS reviewed
        FROM assets a
        LEFT JOIN inventory_records ir
            ON ir.asset_id  = a.id
            AND ir.period_id = $1
        WHERE a.owner_id = $2
			AND ( a.logical_status IN ('active', 'written_off'))
            AND (a.logical_status = 'active' OR ir.id IS NOT NULL)`, periodID, ownerId,
	).Scan(&total, &reviewed)
	return
}

// FindAssetsWithPeriodStatus devuelve todos los activos activos
// con su estado en el período (NULL si no fue revisado aún)
func (r *InventoryRepository) FindAssetsWithPeriodStatus(periodID string, ownerId string) ([]models.AssetInventoryStatus, error) {
	rows, err := r.db.Query(`
        SELECT
            a.id          AS asset_id,
            a.code        AS asset_code,
            a.description AS asset_description,
			a.activation_date AS activation_date,
            ac.name       AS category_name,
            c.name        AS city_name,
            ar.name       AS area_name,
            ir.id         AS record_id,
            ir.confirmed  AS confirmed,
            ir.deactivated AS deactivated,
			ir.has_label AS has_label,
            ir.notes      AS notes,
            ir.recorded_at AS recorded_at
        FROM assets a
        JOIN asset_categories  ac ON ac.id = a.category_id
        JOIN cities            c  ON c.id  = a.city_id
        LEFT JOIN areas        ar ON ar.id = a.area_id
        LEFT JOIN inventory_records ir
            ON ir.asset_id  = a.id
            AND ir.period_id = $1
        WHERE a.owner_id = $2
			AND ( a.logical_status = 'active'
          	 OR (a.logical_status = 'written_off' AND ir.id IS NOT NULL))
        ORDER BY ir.id NULLS LAST, a.code`,
		periodID, ownerId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.AssetInventoryStatus{}
	for rows.Next() {
		var s models.AssetInventoryStatus
		if err := rows.Scan(
			&s.AssetID,
			&s.AssetCode,
			&s.AssetDescription,
			&s.ActivationDate,
			&s.CategoryName,
			&s.CityName,
			&s.AreaName,
			&s.RecordID,
			&s.Confirmed,
			&s.Deactivated,
			&s.HasLabel,
			&s.Notes,
			&s.RecordedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}
