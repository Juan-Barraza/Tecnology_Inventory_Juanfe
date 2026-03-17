package models

import "time"

// AccountingGroup is the parent — one unique code per group.
// Editable from the app (name only).
type AccountingGroup struct {
	ID        int       `db:"id"`
	Code      int64     `db:"code"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
