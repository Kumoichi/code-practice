package persistence

import "code-practice/go/clean-architecture/repository-interface/good/domain"

// CachedBoxRepositoryはdomain.BoxRepositoryをラップしてキャッシュを追加する。
// innerもdomain.BoxRepository型(interface)なので、本物のBoxRepositoryだけでなく
// Mockや他の実装もそのまま包める。CachedBoxRepository自身もFindを持つので、
// domain.BoxRepositoryを満たし、application.NewBoxUseCaseにそのまま渡せる。
type CachedBoxRepository struct {
	inner domain.BoxRepository
	cache map[int]*domain.Box
}

func NewCachedBoxRepository(inner domain.BoxRepository) *CachedBoxRepository {
	return &CachedBoxRepository{inner: inner, cache: map[int]*domain.Box{}}
}

func (r *CachedBoxRepository) Find(id int) (*domain.Box, error) {
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
