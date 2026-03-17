package repository

import (
	"database/sql"

	"inventory-juanfe/models"
)

type AreaRepository struct {
	db *sql.DB
}

func NewAreaRepository(db *sql.DB) *AreaRepository {
	return &AreaRepository{db: db}
}

func (r *AreaRepository) FindAll() ([]models.Area, error) {
	rows, err := r.db.Query(`
        SELECT id, name, description
        FROM areas
        ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var areas []models.Area
	for rows.Next() {
		var a models.Area
		if err := rows.Scan(&a.ID, &a.Name, &a.Description); err != nil {
			return nil, err
		}
		areas = append(areas, a)
	}
	return areas, rows.Err()
}

func (r *AreaRepository) FindByID(id int) (*models.Area, error) {
	var a models.Area
	err := r.db.QueryRow(`
        SELECT id, name, description FROM areas WHERE id = $1`, id,
	).Scan(&a.ID, &a.Name, &a.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}
