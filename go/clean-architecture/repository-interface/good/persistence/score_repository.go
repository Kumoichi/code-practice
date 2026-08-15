// Package persistence is the only place in "good" that is allowed to know
// about PostgreSQL. It implements domain.ScoreRepository.
package persistence

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"code-practice/go/clean-architecture/repository-interface/good/domain"
)

type ScoreRepository struct {
	db *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

// Find implements domain.ScoreRepository.
func (r *ScoreRepository) Find(id int) (*domain.Score, error) {
	row := r.db.QueryRow(`SELECT id, value FROM scores WHERE id = $1`, id)

	var score domain.Score
	if err := row.Scan(&score.ID, &score.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("score not found: id=%d", id)
		}
		return nil, err
	}

	return &score, nil
}
