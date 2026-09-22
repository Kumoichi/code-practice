package persistence

import (
	"testing"

	"code-practice/go/clean-architecture/repository-interface/good/domain"
)

// countingBoxRepositoryはdomain.BoxRepositoryを満たすだけのテスト用フェイク。
// これもDBを使わずに済むのは、CachedBoxRepositoryのinnerがinterface型だから。
type countingBoxRepository struct {
	Box       *domain.Box
	FindCalls int
}

func (r *countingBoxRepository) Find(id int) (*domain.Box, error) {
	r.FindCalls++
	return r.Box, nil
}

func TestCachedBoxRepository_Find_2回目はキャッシュから返り内側は呼ばれない(t *testing.T) {
	inner := &countingBoxRepository{Box: &domain.Box{ID: 1, Number: 5}}
	repository := NewCachedBoxRepository(inner)

	if _, err := repository.Find(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := repository.Find(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inner.FindCalls != 1 {
		t.Errorf("inner.FindCalls = %d, want 1（2回目はキャッシュから返るはず）", inner.FindCalls)
	}
}
