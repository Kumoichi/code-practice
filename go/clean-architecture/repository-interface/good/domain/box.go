// Package domain holds the core business types. It must not import
// anything from persistence — no database/sql, no driver packages.
package domain

// Box is a single box's recorded content.
type Box struct {
	ID     int
	Number int
}
