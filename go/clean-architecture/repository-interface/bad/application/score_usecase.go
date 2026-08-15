package application

import (
	"errors"

	"code-practice/go/clean-architecture/repository-interface/bad/persistence"
)

// ScoreUseCase depends on the concrete *persistence.ScoreRepository type,
// not an interface. It knows this repository is PostgreSQL-backed.
type ScoreUseCase struct {
	repository *persistence.ScoreRepository
}

func NewScoreUseCase(repository *persistence.ScoreRepository) *ScoreUseCase {
	return &ScoreUseCase{repository: repository}
}

func (u *ScoreUseCase) CheckPass(id int) (bool, error) {
	score, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}

	if score.Value < 0 {
		return false, errors.New("score value must not be negative")
	}

	return score.Value >= 100, nil
}
