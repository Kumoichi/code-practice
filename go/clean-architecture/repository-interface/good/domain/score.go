// Package domain holds the core business types. It must not import
// anything from persistence — no database/sql, no driver packages.
package domain

// Score is a single player's typing score.
type Score struct {
	ID    int
	Value int
}
