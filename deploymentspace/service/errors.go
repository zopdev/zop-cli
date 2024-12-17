package service

import "fmt"

type ErrNoItemSelected struct {
	Type string
}

func (e *ErrNoItemSelected) Error() string {
	return fmt.Sprintf("no %s selected", e.Type)
}
