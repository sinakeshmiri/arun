package ports

import "time"

// APIPort is the technology neutral
// port for driving adapters
type APIPort interface {
	AddFunction(name string, source string) error
	RunFunction(name string) (string, error)
	UpdateFunction(name string, duration time.Duration) error
}
