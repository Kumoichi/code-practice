package application

import (
	"errors"
	"testing"

	"code-practice/go/clean-architecture/repository-interface/4-di-wire/domain"
)

// MockBoxRepositoryはdomain.BoxRepositoryを満たしているので、
// NewBoxUseCaseにDIできる
type MockBoxRepository struct {
	Box *domain.Box
	Err error

	FindCalled       bool
	FindCalledWithID int
}

// Findを定義することでMockBoxRepository型はdomain.BoxRepositoryインターフェースを実装しているといえる
func (m *MockBoxRepository) Find(id int) (*domain.Box, error) {
	m.FindCalled = true
	m.FindCalledWithID = id

	if m.Err != nil {
		return nil, m.Err
	}
	return m.Box, nil
}

// PostgreSQLを止めていても、このテストは動く。
func TestBoxUseCase_IsLarge_5は4以上なのでtrue(t *testing.T) {
	mock := &MockBoxRepository{Box: &domain.Box{ID: 1, Number: 5}}
	useCase := NewBoxUseCase(mock) // ← ここでMockをDIしている

	got, err := useCase.IsLarge(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Errorf("IsLarge() = false, want true")
	}
}

func TestBoxUseCase_IsLarge_2は4未満なのでfalse(t *testing.T) {
	mock := &MockBoxRepository{Box: &domain.Box{ID: 2, Number: 2}}
	useCase := NewBoxUseCase(mock)

	got, err := useCase.IsLarge(2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Errorf("IsLarge() = true, want false")
	}
}

func TestBoxUseCase_IsLarge_Repositoryのエラーはそのまま返る(t *testing.T) {
	mock := &MockBoxRepository{Err: errors.New("connection refused")}
	useCase := NewBoxUseCase(mock)

	_, err := useCase.IsLarge(999)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestBoxUseCase_IsLarge_Repositoryが正しいIDで呼ばれる(t *testing.T) {
	mock := &MockBoxRepository{Box: &domain.Box{ID: 42, Number: 4}}
	useCase := NewBoxUseCase(mock)

	if _, err := useCase.IsLarge(42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mock.FindCalled {
		t.Error("Find() was not called")
	}
	if mock.FindCalledWithID != 42 {
		t.Errorf("Find() called with id = %d, want 42", mock.FindCalledWithID)
	}
}
