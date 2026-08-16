package application

import "code-practice/go/clean-architecture/repository-interface/bad/persistence"

// DiceUseCase depends on the concrete *persistence.DiceRepository type,
// not an interface. It knows this repository is PostgreSQL-backed.
type DiceUseCase struct {
	repository *persistence.DiceRepository
}

func NewDiceUseCase(repository *persistence.DiceRepository) *DiceUseCase {
	return &DiceUseCase{repository: repository}
}

// IsBig reports whether the dice roll identified by id came up 4 or higher.
func (u *DiceUseCase) IsBig(id int) (bool, error) {
	roll, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}

	return roll.Pips >= 4, nil
}
