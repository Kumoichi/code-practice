// Package application demonstrates the "no DI" anti-pattern: BoxUseCase
// builds its own dependency instead of receiving one from outside.
// Compare with 3-di-interface/application/box_usecase.go (DI + interface) and
// 2-di-concrete/application/box_usecase.go (DI + concrete type).
package application

import "code-practice/go/clean-architecture/repository-interface/1-no-di/persistence"

type BoxUseCase struct {
	repository *persistence.BoxRepository
}

// DIをしていない: 引数が0個で、persistenceの具体的な関数を名指しで呼んでいる。
func NewBoxUseCase() *BoxUseCase {
	return &BoxUseCase{repository: persistence.NewDefaultBoxRepository()}
}

// IsLarge reports whether the number in the box identified by id is 4 or higher.
func (u *BoxUseCase) IsLarge(id int) (bool, error) {
	box, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}

	return box.Number >= 4, nil
}
