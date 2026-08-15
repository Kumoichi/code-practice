package domain

// ScoreRepository describes what the application needs from a Score store.
// It says nothing about PostgreSQL, SQL, HTTP, or any other infrastructure
// detail — only "a Score can be found by ID".
type ScoreRepository interface {
	Find(id int) (*Score, error)
}
