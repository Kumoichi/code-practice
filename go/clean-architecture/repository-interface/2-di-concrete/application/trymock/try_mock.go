// Package trymock は意図的にコンパイルエラーになります。
//
// 「2-di-concrete の設計ではMockを差し込めない」ことを、文章ではなく
// 実際のコンパイルエラーで示すためのサンプルです。エディタで開くと、
// 下の application.NewBoxUseCase(mock) の行に赤線が出ます。
//
// このファイルだけを独立したパッケージに置いているのは、application
// パッケージに同居させると、そちらのテスト(box_usecase_test.go)まで
// ビルドできなくなり実行不能になってしまうためです。
//
// 詳しい解説は ../mock_cannot_be_injected.txt を参照してください。
package trymock

import (
	"code-practice/go/clean-architecture/repository-interface/2-di-concrete/application"
	"code-practice/go/clean-architecture/repository-interface/2-di-concrete/persistence"
)

// 本物のBoxRepositoryと同じ形のFindを持つMock。
// 3-di-interface 側ではこれと同じ発想のMockがそのまま注入できる。
type MockBoxRepository struct {
	Box *persistence.Box
	Err error
}

func (m *MockBoxRepository) Find(id int) (*persistence.Box, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Box, nil
}

func tryInject() {
	mock := &MockBoxRepository{Box: &persistence.Box{ID: 1, Number: 5}}

	// ここでコンパイルエラー:
	//   cannot use mock (variable of type *MockBoxRepository)
	//   as *persistence.CachedBoxRepository value in argument to application.NewBoxUseCase
	//
	// NewBoxUseCase が interface ではなく *persistence.CachedBoxRepository
	// という具体的な構造体を要求しているため、同じメソッドを持っていても
	// 型が違う時点で受け付けられない。
	_ = application.NewBoxUseCase(mock)
}
