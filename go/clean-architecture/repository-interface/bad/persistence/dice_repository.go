// Package persistence is the concrete, PostgreSQL-specific implementation.
// Unlike "good", there is no domain package here — DiceRoll is defined by
// persistence itself, and application will import this package directly.
package persistence

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type DiceRoll struct {
	ID   int
	Pips int
}

type DiceRepository struct {
	db *sql.DB
}

func NewDiceRepository(db *sql.DB) *DiceRepository {
	return &DiceRepository{db: db}
}

func (r *DiceRepository) Find(id int) (*DiceRoll, error) {
	row := r.db.QueryRow(`SELECT id, pips FROM dice_rolls WHERE id = $1`, id)

	var roll DiceRoll
	if err := row.Scan(&roll.ID, &roll.Pips); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("dice roll not found: id=%d", id)
		}
		return nil, err
	}

	return &roll, nil
}
