// Package domain holds the core business types. It must not import
// anything from persistence — no database/sql, no driver packages.
package domain

// DiceRoll is a single recorded dice roll.
type DiceRoll struct {
	ID   int
	Pips int
}
