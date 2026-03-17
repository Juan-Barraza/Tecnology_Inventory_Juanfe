package repository

import (
	"database/sql"

	"inventory-juanfe/models"
)

type StatusHistoryRepository struct {
	db *sql.DB
}

func NewStatusHistoryRepository(db *sql.DB) *StatusHistoryRepository {
	return &StatusHistoryRepository{db: db}
}

func (r *StatusHistoryRepository) Create(h *models.StatusHistory) error {
	_, err := r.db.Exec(`
        INSERT INTO status_history
            (id, asset_id, previous_status, new_status, notes, recorded_by)
        VALUES ($1, $2, $3, $4, $5, $6)`,
		h.ID, h.AssetID, h.PreviousStatus,
		h.NewStatus, h.Notes, h.RecordedBy,
	)
	return err
}

func (r *StatusHistoryRepository) FindByAssetID(assetID string) ([]models.StatusHistoryDetail, error) {
	rows, err := r.db.Query(`
        SELECT
            sh.id, sh.asset_id, sh.previous_status,
            sh.new_status, sh.notes, sh.recorded_by, sh.created_at,
            u.name AS recorded_by_name
        FROM status_history sh
        JOIN users u ON u.id = sh.recorded_by
        WHERE sh.asset_id = $1
        ORDER BY sh.created_at DESC`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []models.StatusHistoryDetail{}
	for rows.Next() {
		var h models.StatusHistoryDetail
		if err := rows.Scan(
			&h.ID, &h.AssetID, &h.PreviousStatus,
			&h.NewStatus, &h.Notes, &h.RecordedBy, &h.CreatedAt,
			&h.RecordedByName,
		); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, rows.Err()
}
