package application

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"code-practice/go/clean-architecture/repository-interface/bad/persistence"
)

// Unlike good/application/dice_usecase_test.go, this test cannot use a
// mock: NewDiceUseCase only accepts *persistence.DiceRepository, so the
// only way to exercise IsBig is through a real PostgreSQL instance.
// See mock_cannot_be_injected.txt for what happens if you try.
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

func TestDiceUseCase_IsBig_5の目はtrue(t *testing.T) {
	db := openTestDB(t)
	useCase := NewDiceUseCase(persistence.NewDiceRepository(db))

	got, err := useCase.IsBig(1) // seed data: id=1, pips=5

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Errorf("IsBig() = false, want true")
	}
}

func TestDiceUseCase_IsBig_2の目はfalse(t *testing.T) {
	db := openTestDB(t)
	useCase := NewDiceUseCase(persistence.NewDiceRepository(db))

	got, err := useCase.IsBig(2) // seed data: id=2, pips=2

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Errorf("IsBig() = true, want false")
	}
}

func TestDiceUseCase_IsBig_存在しないIDはエラー(t *testing.T) {
	db := openTestDB(t)
	useCase := NewDiceUseCase(persistence.NewDiceRepository(db))

	_, err := useCase.IsBig(9999) // does not exist

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
