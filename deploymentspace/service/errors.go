package service

import "fmt"

// ErrNoItemSelected represents an error that occurs when no item of a specific type is selected.
type ErrNoItemSelected struct {
	Type string
}

// Error returns the error message for ErrNoItemSelected.
//
// This method satisfies the error interface.
//
// Returns:
//   - A formatted error message indicating the type of the unselected item.
func (e *ErrNoItemSelected) Error() string {
	return fmt.Sprintf("no %s selected", e.Type)
}
