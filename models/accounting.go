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

// Para el formulario — grupo con sus subcuentas anidadas
type AssetAccountItem struct {
	ID          int     `json:"id"`
	AccountCode int64   `json:"account_code"`
	OpenLedger  *string `json:"open_ledger"`
}

type AccountingGroupWithAccounts struct {
	ID       int                `json:"id"`
	Code     int64              `json:"code"`
	Name     string             `json:"name"`
	Accounts []AssetAccountItem `json:"accounts"`
}
