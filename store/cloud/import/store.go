package _import

import (
	"database/sql"
	"gofr.dev/pkg/gofr"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetAccounts(ctx *gofr.Context) ([]string, error) {
	ans := make([]string, 0)

	rows, err := s.db.Query("SELECT account_id FROM credentials")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var account_id string
		err = rows.Scan(&account_id)
		if err != nil {
			return nil, err
		}
		ans = append(ans, account_id)
	}

	return ans, nil
}
