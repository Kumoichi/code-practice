package application

import (
	"errors"
	"testing"

	"code-practice/go/clean-architecture/repository-interface/good/domain"
)

// MockDiceRepositoryはdomain.DiceRepositoryを満たしているので、
// NewDiceUseCaseにDIできる
type MockDiceRepository struct {
	Roll *domain.DiceRoll
	Err  error

	FindCalled       bool
	FindCalledWithID int
}

// Findを定義することでMockDiceRepository型はdoamin.DiceRollインターフェースを実装しているといえる
func (m *MockDiceRepository) Find(id int) (*domain.DiceRoll, error) {
	m.FindCalled = true
	m.FindCalledWithID = id

	if m.Err != nil {
		return nil, m.Err
	}
	return m.Roll, nil
}

// PostgreSQLを止めていても、このテストは動く。
func TestDiceUseCase_IsBig_5の目はtrue(t *testing.T) {
	mock := &MockDiceRepository{Roll: &domain.DiceRoll{ID: 1, Pips: 5}}
	useCase := NewDiceUseCase(mock) // ← ここでMockをDIしている

	got, err := useCase.IsBig(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Errorf("IsBig() = false, want true")
	}
}

func TestDiceUseCase_IsBig_2の目はfalse(t *testing.T) {
	mock := &MockDiceRepository{Roll: &domain.DiceRoll{ID: 2, Pips: 2}}
	useCase := NewDiceUseCase(mock)

	got, err := useCase.IsBig(2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Errorf("IsBig() = true, want false")
	}
}

func TestDiceUseCase_IsBig_Repositoryのエラーはそのまま返る(t *testing.T) {
	mock := &MockDiceRepository{Err: errors.New("connection refused")}
	useCase := NewDiceUseCase(mock)

	_, err := useCase.IsBig(999)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestDiceUseCase_IsBig_Repositoryが正しいIDで呼ばれる(t *testing.T) {
	mock := &MockDiceRepository{Roll: &domain.DiceRoll{ID: 42, Pips: 4}}
	useCase := NewDiceUseCase(mock)

	if _, err := useCase.IsBig(42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mock.FindCalled {
		t.Error("Find() was not called")
	}
	if mock.FindCalledWithID != 42 {
		t.Errorf("Find() called with id = %d, want 42", mock.FindCalledWithID)
	}
}
