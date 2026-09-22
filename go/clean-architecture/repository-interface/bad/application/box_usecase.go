package application

import "code-practice/go/clean-architecture/repository-interface/bad/persistence"

// BoxUseCase depends on the concrete *persistence.BoxRepository type,
// not an interface. It knows this repository is PostgreSQL-backed.
type BoxUseCase struct {
	repository *persistence.BoxRepository
}

func NewBoxUseCase(repository *persistence.BoxRepository) *BoxUseCase {
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
