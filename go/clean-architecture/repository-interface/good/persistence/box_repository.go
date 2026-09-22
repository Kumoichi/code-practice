// Package persistence is the only place in "good" that is allowed to know
// about PostgreSQL. It implements domain.BoxRepository.
package persistence

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"code-practice/go/clean-architecture/repository-interface/good/domain"
)

type BoxRepository struct {
	db *sql.DB
}

func NewBoxRepository(db *sql.DB) *BoxRepository {
	return &BoxRepository{db: db}
}

// Find implements domain.BoxRepository.
func (r *BoxRepository) Find(id int) (*domain.Box, error) {
	row := r.db.QueryRow(`SELECT id, number FROM boxes WHERE id = $1`, id)

	var box domain.Box
	if err := row.Scan(&box.ID, &box.Number); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("box not found: id=%d", id)
		}
		return nil, err
	}

	return &box, nil
}
