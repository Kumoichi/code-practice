package main

import (
	"fmt"

	"code-practice/go/clean-architecture/repository-interface/nodi/application"
)

func main() {
	useCase := application.NewBoxUseCase()

	for _, id := range []int{1, 2, 3} {
		large, err := useCase.IsLarge(id)
		if err != nil {
			fmt.Printf("id=%d error=%v\n", id, err)
			continue
		}
		fmt.Printf("id=%d large=%v\n", id, large)
	}
}
