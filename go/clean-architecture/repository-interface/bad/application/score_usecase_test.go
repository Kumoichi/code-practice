package application

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"code-practice/go/clean-architecture/repository-interface/bad/persistence"
)

// Unlike good/application/score_usecase_test.go, this test cannot use a
// mock: NewScoreUseCase only accepts *persistence.ScoreRepository, so the
// only way to exercise CheckPass is through a real PostgreSQL instance.
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

func TestScoreUseCase_CheckPass_120点はtrue(t *testing.T) {
	db := openTestDB(t)
	useCase := NewScoreUseCase(persistence.NewScoreRepository(db))

	passed, err := useCase.CheckPass(1) // seed data: id=1, value=120

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !passed {
		t.Errorf("passed = false, want true")
	}
}

func TestScoreUseCase_CheckPass_80点はfalse(t *testing.T) {
	db := openTestDB(t)
	useCase := NewScoreUseCase(persistence.NewScoreRepository(db))

	passed, err := useCase.CheckPass(2) // seed data: id=2, value=80

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if passed {
		t.Errorf("passed = true, want false")
	}
}

func TestScoreUseCase_CheckPass_マイナス点はerror(t *testing.T) {
	db := openTestDB(t)
	useCase := NewScoreUseCase(persistence.NewScoreRepository(db))

	_, err := useCase.CheckPass(3) // seed data: id=3, value=-10

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestScoreUseCase_CheckPass_Repositoryのエラーはそのまま返る(t *testing.T) {
	db := openTestDB(t)
	useCase := NewScoreUseCase(persistence.NewScoreRepository(db))

	_, err := useCase.CheckPass(9999) // does not exist

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
