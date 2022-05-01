package wlinterpreter

import "errors"

var ErrNoValue = errors.New(`value doesn't exists`)
var ErrNotAllocatedBlock = errors.New(`block did not allocated`)

// Storage is interface for storing additional data for interpreter. Default implied it implement FIFO interface.
type Storage interface {
	// Alloc new block in storage
	Alloc(name string) error
	// Append adds value to the end of block
	Append(name string, val int) error
	// Pop removes and returns the object at the end of block in storage
	Pop(name string) (int, error)
	// Last method returns either value of last element if it exists or ErrNoValue if not
	Last(name string) (int, error)
	// Empty checks does storage block exists and does is has any value in block
	Empty(s string) bool
	// Update last value
	Update(name string, val int) error
}
