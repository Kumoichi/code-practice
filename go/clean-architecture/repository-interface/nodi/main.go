package main

import (
	"fmt"

	"code-practice/go/clean-architecture/repository-interface/nodi/application"
)

func main() {
	useCase := application.NewDiceUseCase()

	for _, id := range []int{1, 2, 3} {
		big, err := useCase.IsBig(id)
		if err != nil {
			fmt.Printf("id=%d error=%v\n", id, err)
			continue
		}
		fmt.Printf("id=%d big=%v\n", id, big)
	}
}
