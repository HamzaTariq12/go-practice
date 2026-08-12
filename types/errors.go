package types

import "fmt"

type PathError struct {
	Op   string // what operation failed, e.g. "open"
	Path string // which path
	Err  error  // underlying cause
}

func (e *PathError) Error() string {
	return fmt.Sprintf("%s %s: %v", e.Op, e.Path, e.Err)
}

func (e *PathError) Unwrap() error {
	return e.Err
}
