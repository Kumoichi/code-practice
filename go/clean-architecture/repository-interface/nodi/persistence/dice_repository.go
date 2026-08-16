package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DiceRoll struct {
	ID   int
	Pips int
}

type DiceRepository struct {
	db *sql.DB
}

// DIをしている: dbを外から受け取るだけで、自分では作らない。
func NewDiceRepository(db *sql.DB) *DiceRepository {
	return &DiceRepository{db: db}
}

// DIをしていない: 引数が0個で、dsnの組み立てからsql.Openまで自分でやる。
func NewDefaultDiceRepository() *DiceRepository {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/code_practice?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return NewDiceRepository(db)
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
