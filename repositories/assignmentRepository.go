package repository

import (
	"database/sql"
	"time"

	"inventory-juanfe/models"
)

type AssignmentRepository struct {
	db *sql.DB
}

func NewAssignmentRepository(db *sql.DB) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

const assignmentSelectBase = `
    SELECT
        a.id, a.asset_id, a.responsible_name, a.responsible_position,
        a.assigned_at, a.deactivated_at, a.deactivation_reason,
        a.status, a.created_by, a.created_at,
        ast.code        AS asset_code,
        ast.description AS asset_description,
        u.name          AS created_by_name
    FROM assignments a
    JOIN assets ast ON ast.id = a.asset_id
    JOIN users  u   ON u.id   = a.created_by`

func (r *AssignmentRepository) FindByAssetID(assetID string) ([]models.AssignmentDetail, error) {
	rows, err := r.db.Query(
		assignmentSelectBase+` WHERE a.asset_id = $1 ORDER BY a.created_at DESC`,
		assetID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAssignments(rows)
}

func (r *AssignmentRepository) FindActiveByAssetID(assetID string) (*models.AssignmentDetail, error) {
	row := r.db.QueryRow(
		assignmentSelectBase+` WHERE a.asset_id = $1 AND a.status = 'active'`,
		assetID,
	)
	var a models.AssignmentDetail
	if err := scanAssignment(row, &a); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AssignmentRepository) Create(a *models.Assignment) error {
	_, err := r.db.Exec(`
        INSERT INTO assignments
            (id, asset_id, responsible_name, responsible_position,
             assigned_at, status, created_by)
        VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		a.ID, a.AssetID, a.ResponsibleName, a.ResponsiblePosition,
		a.AssignedAt, a.Status, a.CreatedBy,
	)
	return err
}

func (r *AssignmentRepository) Release(id string, deactivatedAt time.Time, reason *string) error {
	_, err := r.db.Exec(`
        UPDATE assignments
        SET status = 'released', deactivated_at = $1, deactivation_reason = $2
        WHERE id = $3`,
		deactivatedAt, reason, id,
	)
	return err
}

func (r *AssignmentRepository) WriteOff(assetID string, deactivatedAt time.Time, reason *string) error {
	_, err := r.db.Exec(`
        UPDATE assignments
        SET status = 'written_off', deactivated_at = $1, deactivation_reason = $2
        WHERE asset_id = $3 AND status = 'active'`,
		deactivatedAt, reason, assetID,
	)
	return err
}

// ── scan helpers ──────────────────────────────────────────────

type assignmentScanner interface {
	Scan(dest ...any) error
}

func scanAssignment(s assignmentScanner, a *models.AssignmentDetail) error {
	return s.Scan(
		&a.ID, &a.AssetID, &a.ResponsibleName, &a.ResponsiblePosition,
		&a.AssignedAt, &a.DeactivatedAt, &a.DeactivationReason,
		&a.Status, &a.CreatedBy, &a.CreatedAt,
		&a.AssetCode, &a.AssetDescription, &a.CreatedByName,
	)
}

func scanAssignments(rows *sql.Rows) ([]models.AssignmentDetail, error) {
	result := []models.AssignmentDetail{}
	for rows.Next() {
		var a models.AssignmentDetail
		if err := scanAssignment(rows, &a); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}
