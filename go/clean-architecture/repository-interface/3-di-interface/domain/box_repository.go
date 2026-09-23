package domain

// BoxRepository describes what the application needs from a box
// store. It says nothing about PostgreSQL, SQL, HTTP, or any other
// infrastructure detail — only "a Box can be found by ID".
type BoxRepository interface {
	Find(id int) (*Box, error)
}
