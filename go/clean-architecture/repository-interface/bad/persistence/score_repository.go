// Package persistence is the concrete, PostgreSQL-specific implementation.
// Unlike "good", there is no domain package here — Score is defined by
// persistence itself, and application will import this package directly.
package persistence

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type Score struct {
	ID    int
	Value int
}

type ScoreRepository struct {
	db *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

func (r *ScoreRepository) Find(id int) (*Score, error) {
	row := r.db.QueryRow(`SELECT id, value FROM scores WHERE id = $1`, id)

	var score Score
	if err := row.Scan(&score.ID, &score.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("score not found: id=%d", id)
		}
		return nil, err
	}

	return &score, nil
}
