package repository

import (
	"database/sql"

	"inventory-juanfe/models"
)

type CityRepository struct {
	db *sql.DB
}

func NewCityRepository(db *sql.DB) *CityRepository {
	return &CityRepository{db: db}
}

func (r *CityRepository) FindAll() ([]models.City, error) {
	rows, err := r.db.Query(`
        SELECT id, name, department
        FROM cities
        ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var c models.City
		if err := rows.Scan(&c.ID, &c.Name, &c.Department); err != nil {
			return nil, err
		}
		cities = append(cities, c)
	}
	return cities, rows.Err()
}

func (r *CityRepository) FindByID(id int) (*models.City, error) {
	var c models.City
	err := r.db.QueryRow(`
        SELECT id, name, department FROM cities WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Department)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}
