// Package persistence is the only place in "good" that is allowed to know
// about PostgreSQL. It implements domain.DiceRepository.
package persistence

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"code-practice/go/clean-architecture/repository-interface/good/domain"
)

type DiceRepository struct {
	db *sql.DB
}

func NewDiceRepository(db *sql.DB) *DiceRepository {
	return &DiceRepository{db: db}
}

// Find implements domain.DiceRepository.
func (r *DiceRepository) Find(id int) (*domain.DiceRoll, error) {
	row := r.db.QueryRow(`SELECT id, pips FROM dice_rolls WHERE id = $1`, id)

	var roll domain.DiceRoll
	if err := row.Scan(&roll.ID, &roll.Pips); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("dice roll not found: id=%d", id)
		}
		return nil, err
	}

	return &roll, nil
}
