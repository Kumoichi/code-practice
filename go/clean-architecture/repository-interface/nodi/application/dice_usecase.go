// Package application demonstrates the "no DI" anti-pattern: DiceUseCase
// builds its own dependency instead of receiving one from outside.
// Compare with good/application/dice_usecase.go (DI + interface) and
// bad/application/dice_usecase.go (DI + concrete type).
package application

import "code-practice/go/clean-architecture/repository-interface/nodi/persistence"

type DiceUseCase struct {
	repository *persistence.DiceRepository
}

// DIをしていない: 引数が0個で、persistenceの具体的な関数を名指しで呼んでいる。
func NewDiceUseCase() *DiceUseCase {
	return &DiceUseCase{repository: persistence.NewDefaultDiceRepository()}
}

// IsBig reports whether the dice roll identified by id came up 4 or higher.
func (u *DiceUseCase) IsBig(id int) (bool, error) {
	roll, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}

	return roll.Pips >= 4, nil
}
