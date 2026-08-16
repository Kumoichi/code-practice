package application

import "code-practice/go/clean-architecture/repository-interface/bad/persistence"

type MockDiceRepository struct {
	Roll *persistence.DiceRoll
	Err  error
}

func (m *MockDiceRepository) Find(id int) (*persistence.DiceRoll, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Roll, nil
}

func tryInject() {
	mock := &MockDiceRepository{Roll: &persistence.DiceRoll{ID: 1, Pips: 5}}
	_ = NewDiceUseCase(mock)
}
