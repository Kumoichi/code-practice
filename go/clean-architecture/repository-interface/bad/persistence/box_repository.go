// Package persistence is the concrete, PostgreSQL-specific implementation.
// Unlike "good", there is no domain package here — Box is defined by
// persistence itself, and application will import this package directly.
package persistence

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type Box struct {
	ID     int
	Number int
}

type BoxRepository struct {
	db *sql.DB
}

func NewBoxRepository(db *sql.DB) *BoxRepository {
	return &BoxRepository{db: db}
}

func (r *BoxRepository) Find(id int) (*Box, error) {
	row := r.db.QueryRow(`SELECT id, number FROM boxes WHERE id = $1`, id)

	var box Box
	if err := row.Scan(&box.ID, &box.Number); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("box not found: id=%d", id)
		}
		return nil, err
	}

	return &box, nil
}
