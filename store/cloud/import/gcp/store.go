// Package gcp implements the AccountGetter interface to fetch all accounts from the gcloud database
// at ~/.config/gcloud/credentials.db, reading the account_id and value from the database
package gcp

import (
	"database/sql"

	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/models"
)

// Store is a struct that holds the database connection and implements the AccountGetter interface
type Store struct {
	db *sql.DB
}

// New creates a new Store struct to fetch all accounts from the gcloud database
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// GetAccounts returns a list of accounts from the database reading the account_id and value
// from the gcloud database at ~/.config/gcloud/credentials.db
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
