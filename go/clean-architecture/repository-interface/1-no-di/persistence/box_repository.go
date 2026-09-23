package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Box struct {
	ID     int
	Number int
}

type BoxRepository struct {
	db *sql.DB
}

// DIをしている: dbを外から受け取るだけで、自分では作らない。
func NewBoxRepository(db *sql.DB) *BoxRepository {
	return &BoxRepository{db: db}
}

// DIをしていない: 引数が0個で、dsnの組み立てからsql.Openまで自分でやる。
func NewDefaultBoxRepository() *BoxRepository {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/code_practice?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return NewBoxRepository(db)
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
