package persistence

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// This test hits real PostgreSQL. It requires `docker compose up -d`
// to have been run from the repository root. Compare with
// good/application/dice_usecase_test.go, which needs no database at all.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/code_practice?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := db.Ping(); err != nil {
		t.Fatalf("postgres is required for this test (run `docker compose up -d` first): %v", err)
	}

	return db
}

func TestDiceRepository_Find_存在する目を取得できる(t *testing.T) {
	db := openTestDB(t)
	repository := NewDiceRepository(db)

	roll, err := repository.Find(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if roll.ID != 1 || roll.Pips != 5 {
		t.Errorf("got %+v, want {ID:1 Pips:5}", roll)
	}
}

func TestDiceRepository_Find_存在しないIDはエラー(t *testing.T) {
	db := openTestDB(t)
	repository := NewDiceRepository(db)

	_, err := repository.Find(9999)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
