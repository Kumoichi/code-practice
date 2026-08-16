package application

import "code-practice/go/clean-architecture/repository-interface/good/domain"

// repositoryはinterfaceなので、DiceUseCaseは中身が本物のDBなのかMockなのかを知らない
type DiceUseCase struct {
	repository domain.DiceRepository
}

// repositoryを自分でnewせず、外から受け取る(=DI)。呼び出す側が
// 本物のpersistence.DiceRepositoryを渡すかMockDiceRepositoryを渡すかを決める
func NewDiceUseCase(repository domain.DiceRepository) *DiceUseCase {
	return &DiceUseCase{repository: repository}
}

func (u *DiceUseCase) IsBig(id int) (bool, error) {
	// u.repositoryの中身は実行時にしか決まらない。NewDiceUseCaseで注入された
	// 具体型(本物 or Mock)のFindが動的に呼ばれる
	roll, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}

	return roll.Pips >= 4, nil
}
