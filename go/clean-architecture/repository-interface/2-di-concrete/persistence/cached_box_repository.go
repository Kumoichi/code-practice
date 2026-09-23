package persistence

// CachedBoxRepositoryはBoxRepositoryをラップしてキャッシュを追加する。
// innerがinterfaceではなく*BoxRepositoryという具体型なので、テスト用のフェイクを
// 包むことはできず、本物のBoxRepository(=本物のDB接続)しか渡せない。
// 比較: 3-di-interface/persistence/cached_box_repository.go の inner は
// domain.BoxRepository(interface)なので、Mockでも何でも包める。
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
