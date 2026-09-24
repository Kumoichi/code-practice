//go:build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"

	"code-practice/go/clean-architecture/repository-interface/4-di-wire/application"
	"code-practice/go/clean-architecture/repository-interface/4-di-wire/domain"
	"code-practice/go/clean-architecture/repository-interface/4-di-wire/persistence"
)

// newCachedRepository は素のBoxRepositoryをキャッシュで包み、domain.BoxRepositoryとして返す
func newCachedRepository(inner *persistence.BoxRepository) domain.BoxRepository {
	return persistence.NewCachedBoxRepository(inner)
}

// InitializeBoxUseCase は「BoxUseCaseが欲しい。材料はこれ」と宣言するだけのinjector。
// 中身の組み立ては `wire` コマンドが wire_gen.go に生成する
func InitializeBoxUseCase(db *sql.DB) *application.BoxUseCase {
	wire.Build(
		persistence.NewBoxRepository,
		newCachedRepository,
		application.NewBoxUseCase,
	)
	return nil
}
