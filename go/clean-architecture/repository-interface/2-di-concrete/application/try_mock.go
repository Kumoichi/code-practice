package application

import "code-practice/go/clean-architecture/repository-interface/2-di-concrete/persistence"

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
	_ = NewBoxUseCase(mock)
}
