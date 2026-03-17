package repository

import (
	"database/sql"

	"inventory-juanfe/models"
)

type AccountingGroupRepository struct {
	db *sql.DB
}

func NewAccountingGroupRepository(db *sql.DB) *AccountingGroupRepository {
	return &AccountingGroupRepository{db: db}
}

func (r *AccountingGroupRepository) FindAll() ([]models.AccountingGroup, error) {
	rows, err := r.db.Query(`
        SELECT id, code, name, created_at, updated_at
        FROM accounting_groups
        ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.AccountingGroup
	for rows.Next() {
		var g models.AccountingGroup
		if err := rows.Scan(
			&g.ID, &g.Code, &g.Name, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

// FindAllWithAccounts devuelve los grupos con sus subcuentas
// para el formulario de creación de activos
func (r *AccountingGroupRepository) FindAllWithAccounts() ([]models.AccountingGroupWithAccounts, error) {
	rows, err := r.db.Query(`
        SELECT
            ag.id, ag.code, ag.name,
            aa.id           AS account_id,
            aa.account_code,
            aa.open_ledger
        FROM accounting_groups ag
        JOIN asset_accounts aa ON aa.accounting_group_id = ag.id
        ORDER BY ag.code, aa.account_code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grouped := make(map[int]*models.AccountingGroupWithAccounts)
	var order []int

	for rows.Next() {
		var (
			gID, aID     int
			gCode, aCode int64
			gName        string
			openLedger   *string
		)
		if err := rows.Scan(
			&gID, &gCode, &gName,
			&aID, &aCode, &openLedger,
		); err != nil {
			return nil, err
		}

		if _, exists := grouped[gID]; !exists {
			grouped[gID] = &models.AccountingGroupWithAccounts{
				ID:   gID,
				Code: gCode,
				Name: gName,
			}
			order = append(order, gID)
		}

		grouped[gID].Accounts = append(grouped[gID].Accounts, models.AssetAccountItem{
			ID:          aID,
			AccountCode: aCode,
			OpenLedger:  openLedger,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]models.AccountingGroupWithAccounts, 0, len(order))
	for _, id := range order {
		result = append(result, *grouped[id])
	}
	return result, nil
}

func (r *AccountingGroupRepository) UpdateName(id int, name string) error {
	_, err := r.db.Exec(`
        UPDATE accounting_groups SET name = $1 WHERE id = $2`,
		name, id,
	)
	return err
}
