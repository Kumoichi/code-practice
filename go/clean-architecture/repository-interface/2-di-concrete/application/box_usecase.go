package application

import "code-practice/go/clean-architecture/repository-interface/2-di-concrete/persistence"

// BoxUseCase depends on the concrete *persistence.CachedBoxRepository type,
// not an interface. It knows this repository is PostgreSQL-backed and cached.
//
// キャッシュを追加するために、フィールドの型とコンストラクタの引数の型を
// *persistence.BoxRepository → *persistence.CachedBoxRepository へ
// 書き換える必要があった。DIはしているのに、application層のコードに
// 手を入れる羽目になっている。
// 比較: 3-di-interface は domain.BoxRepository(interface)のままなので、
// このファイルを1行も変更せずにキャッシュを差し込めた。
type BoxUseCase struct {
	repository *persistence.CachedBoxRepository
}

func NewBoxUseCase(repository *persistence.CachedBoxRepository) *BoxUseCase {
	return &BoxUseCase{repository: repository}
}

// IsLarge reports whether the number in the box identified by id is 4 or higher.
func (u *BoxUseCase) IsLarge(id int) (bool, error) {
	box, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}

	return box.Number >= 4, nil
}
