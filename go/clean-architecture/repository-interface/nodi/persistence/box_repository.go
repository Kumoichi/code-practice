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

// CachedBoxRepositoryはBoxRepositoryをラップしてキャッシュを追加する。
// innerがinterfaceではなく*BoxRepositoryという具体型なので、テスト用のフェイクを
// 包むことはできず、本物のBoxRepository(=本物のDB接続)しか渡せない。
type CachedBoxRepository struct {
	inner *BoxRepository
	cache map[int]*Box
}

func NewCachedBoxRepository(inner *BoxRepository) *CachedBoxRepository {
	return &CachedBoxRepository{inner: inner, cache: map[int]*Box{}}
}

func (r *CachedBoxRepository) Find(id int) (*Box, error) {
	if box, ok := r.cache[id]; ok {
		return box, nil
	}

	box, err := r.inner.Find(id)
	if err != nil {
		return nil, err
	}

	r.cache[id] = box
	return box, nil
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
