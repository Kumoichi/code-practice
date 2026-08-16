package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"code-practice/go/clean-architecture/repository-interface/bad/application"
	"code-practice/go/clean-architecture/repository-interface/bad/persistence"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/code_practice?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	repository := persistence.NewDiceRepository(db)
	useCase := application.NewDiceUseCase(repository)

	for _, id := range []int{1, 2, 3} {
		big, err := useCase.IsBig(id)
		if err != nil {
			fmt.Printf("id=%d error=%v\n", id, err)
			continue
		}
		fmt.Printf("id=%d big=%v\n", id, big)
	}
}
