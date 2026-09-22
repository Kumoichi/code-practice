package application

import "code-practice/go/clean-architecture/repository-interface/good/domain"

// repositoryはinterfaceなので、BoxUseCaseは中身が本物のDBなのかMockなのかを知らない
type BoxUseCase struct {
	repository domain.BoxRepository
}

// repositoryを自分でnewせず、外から受け取る(=DI)。呼び出す側が
// 本物のpersistence.BoxRepositoryを渡すかMockBoxRepositoryを渡すかを決める
func NewBoxUseCase(repository domain.BoxRepository) *BoxUseCase {
	return &BoxUseCase{repository: repository}
}

// IsLarge reports whether the number in the box identified by id is 4 or higher.
func (u *BoxUseCase) IsLarge(id int) (bool, error) {
	// u.repositoryの中身は実行時にしか決まらない。NewBoxUseCaseで注入された
	// 具体型(本物 or Mock)のFindが動的に呼ばれる
	box, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}

	return box.Number >= 4, nil
}
