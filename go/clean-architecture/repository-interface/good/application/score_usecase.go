package application

import (
	"errors"

	"code-practice/go/clean-architecture/repository-interface/good/domain"
)

// ScoreUseCase depends only on the domain.ScoreRepository interface.
// It has no idea whether scores come from PostgreSQL, memory, or a mock.
type ScoreUseCase struct {
	repository domain.ScoreRepository
}

func NewScoreUseCase(repository domain.ScoreRepository) *ScoreUseCase {
	return &ScoreUseCase{repository: repository}
}

// CheckPass reports whether the score identified by id is a passing score
// (>= 100). A negative score value is treated as invalid data.
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
