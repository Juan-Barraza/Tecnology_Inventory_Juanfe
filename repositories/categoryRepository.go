package repository

import (
	"database/sql"

	"inventory-juanfe/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]models.AssetCategory, error) {
	rows, err := r.db.Query(`
        SELECT id, name
        FROM asset_categories
        ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.AssetCategory
	for rows.Next() {
		var c models.AssetCategory
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *CategoryRepository) FindByID(id int) (*models.AssetCategory, error) {
	var c models.AssetCategory
	err := r.db.QueryRow(`
        SELECT id, name FROM asset_categories WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}
