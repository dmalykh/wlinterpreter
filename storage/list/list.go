package list

import (
	"errors"
	"github.com/dmalykh/wlinterpreter"
)

var ErrNothingToRemove = errors.New(`nothing to remove, block already empty`)

type element struct {
	val  int32
	prev *element
}

func New() wlinterpreter.Storage {
	return &ListStorage{
		data: make(map[string]*element),
	}
}

// ListStorage implements Storage interface. Uses single linked list for storing data.
type ListStorage struct {
	data map[string]*element
}

func (l *ListStorage) Empty(s string) bool {
	// Check explicitly for better code reading
	if _, exists := l.data[s]; exists {
		return false
	}
	return l.data[s] == nil
}

func (l *ListStorage) Alloc(name string) error {
	l.data[name] = nil
	return nil
}

func (l *ListStorage) Append(name string, val int32) error {
	if _, exists := l.data[name]; !exists {
		return wlinterpreter.ErrNotAllocatedBlock
	}
	var el element
	el.val = val
	if prev := l.data[name]; prev != nil {
		el.prev = prev
	}
	l.data[name] = &el
	return nil
}

func (l *ListStorage) Pop(name string) (int32, error) {
	if _, exists := l.data[name]; !exists {
		return 0, wlinterpreter.ErrNotAllocatedBlock
	}
	if l.data[name] == nil {
		return 0, wlinterpreter.ErrNoValue
	}
	var val = l.data[name].val

	return val, l.remove(l.data[name])

}

// Remove element. If element doesn't exists, returns error, otherwise change current element for its prev element.
func (l *ListStorage) remove(el *element) error {
	if el == nil {
		return ErrNothingToRemove
	}
	el = el.prev
	return nil
}

// Last returns value of stored list
func (l *ListStorage) Last(name string) (int32, error) {
	if _, exists := l.data[name]; !exists {
		return 0, wlinterpreter.ErrNotAllocatedBlock
	}
	if l.data[name] == nil {
		return 0, wlinterpreter.ErrNoValue
	}
	return l.data[name].val, nil
}
