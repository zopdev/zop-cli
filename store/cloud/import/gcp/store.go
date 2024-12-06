package gcp

import (
	"database/sql"

	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/models"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetAccounts(ctx *gofr.Context) ([]models.AccountStore, error) {
	ans := make([]models.AccountStore, 0)

	rows, err := s.db.Query("SELECT account_id, value FROM credentials")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var acc models.AccountStore

		err = rows.Scan(&acc.AccountID, &acc.Value)
		if err != nil {
			return nil, err
		}

		ans = append(ans, acc)
	}

	return ans, nil
}
