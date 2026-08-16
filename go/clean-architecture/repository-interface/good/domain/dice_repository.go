package domain

// DiceRepository describes what the application needs from a dice roll
// store. It says nothing about PostgreSQL, SQL, HTTP, or any other
// infrastructure detail — only "a DiceRoll can be found by ID".
type DiceRepository interface {
	Find(id int) (*DiceRoll, error)
}
