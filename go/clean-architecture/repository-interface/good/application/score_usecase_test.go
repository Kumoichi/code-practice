package application

import (
	"errors"
	"testing"

	"code-practice/go/clean-architecture/repository-interface/good/domain"
)

// MockScoreRepository is a hand-written test double. It satisfies
// domain.ScoreRepository, so it can be passed to NewScoreUseCase exactly
// like the real PostgreSQL-backed repository — that's the whole point.
type MockScoreRepository struct {
	Score *domain.Score
	Err   error

	FindCalled       bool
	FindCalledWithID int
}

func (m *MockScoreRepository) Find(id int) (*domain.Score, error) {
	m.FindCalled = true
	m.FindCalledWithID = id

	if m.Err != nil {
		return nil, m.Err
	}
	return m.Score, nil
}

// No database, no docker, no network — this test runs even with
// PostgreSQL stopped. Compare with bad/application/score_usecase_test.go.

func TestScoreUseCase_CheckPass_120点はtrue(t *testing.T) {
	mock := &MockScoreRepository{Score: &domain.Score{ID: 1, Value: 120}}
	useCase := NewScoreUseCase(mock)

	passed, err := useCase.CheckPass(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !passed {
		t.Errorf("passed = false, want true")
	}
}

func TestScoreUseCase_CheckPass_80点はfalse(t *testing.T) {
	mock := &MockScoreRepository{Score: &domain.Score{ID: 2, Value: 80}}
	useCase := NewScoreUseCase(mock)

	passed, err := useCase.CheckPass(2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if passed {
		t.Errorf("passed = true, want false")
	}
}

func TestScoreUseCase_CheckPass_マイナス点はerror(t *testing.T) {
	mock := &MockScoreRepository{Score: &domain.Score{ID: 3, Value: -10}}
	useCase := NewScoreUseCase(mock)

	_, err := useCase.CheckPass(3)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestScoreUseCase_CheckPass_Repositoryのエラーはそのまま返る(t *testing.T) {
	mock := &MockScoreRepository{Err: errors.New("connection refused")}
	useCase := NewScoreUseCase(mock)

	_, err := useCase.CheckPass(999)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestScoreUseCase_CheckPass_Repositoryが正しいIDで呼ばれる(t *testing.T) {
	mock := &MockScoreRepository{Score: &domain.Score{ID: 42, Value: 100}}
	useCase := NewScoreUseCase(mock)

	if _, err := useCase.CheckPass(42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mock.FindCalled {
		t.Error("Find() was not called")
	}
	if mock.FindCalledWithID != 42 {
		t.Errorf("Find() called with id = %d, want 42", mock.FindCalledWithID)
	}
}
