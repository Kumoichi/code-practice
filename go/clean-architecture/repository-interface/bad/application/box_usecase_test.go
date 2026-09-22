package application

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"code-practice/go/clean-architecture/repository-interface/bad/persistence"
)

// Unlike good/application/box_usecase_test.go, this test cannot use a
// mock: NewBoxUseCase only accepts *persistence.BoxRepository, so the
// only way to exercise IsLarge is through a real PostgreSQL instance.
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

func TestBoxUseCase_IsLarge_5は4以上なのでtrue(t *testing.T) {
	db := openTestDB(t)
	useCase := NewBoxUseCase(persistence.NewBoxRepository(db))

	got, err := useCase.IsLarge(1) // seed data: id=1, number=5

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Errorf("IsLarge() = false, want true")
	}
}

func TestBoxUseCase_IsLarge_2は4未満なのでfalse(t *testing.T) {
	db := openTestDB(t)
	useCase := NewBoxUseCase(persistence.NewBoxRepository(db))

	got, err := useCase.IsLarge(2) // seed data: id=2, number=2

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Errorf("IsLarge() = true, want false")
	}
}

func TestBoxUseCase_IsLarge_存在しないIDはエラー(t *testing.T) {
	db := openTestDB(t)
	useCase := NewBoxUseCase(persistence.NewBoxRepository(db))

	_, err := useCase.IsLarge(9999) // does not exist

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
