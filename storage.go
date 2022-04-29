package wlinterpreter

import "errors"

var ErrNoValue = errors.New(`value doesn't exists`)
var ErrNotAllocatedBlock = errors.New(`block did not allocated`)

// Storage is interface for storing additional data for interpreter
type Storage interface {
	// Alloc new block in storage
	Alloc(name string) error
	// Append value to the end of block
	Append(name string, val int32) error
	// Pop removes and returns the object at the end of block in storage
	Pop(name string) (int32, error)
	// Last method returns either value of last element if it exist or ErrNoValue if not
	Last(name string) (int32, error)
	// Empty checks does storage block exists and does is has any value in block
	Empty(s string) bool
}
