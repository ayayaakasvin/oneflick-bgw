package worker

import (
	"errors"
	"fmt"
)

type FatalError struct {
	msg string
}

func (fe *FatalError) Error() string {
	return fmt.Sprintf("fatal error: %s", fe.msg)
}

func NewFatalError(msg string) error {
	return &FatalError{msg}
}

// Added sentinel and Is implementation so errors.Is works.
var ErrFatal = errors.New("fatal error")

func (fe *FatalError) Is(target error) bool {
	if target == ErrFatal {
		return true
	}
	switch target.(type) {
	case *FatalError:
		return true
	default:
		return false
	}
}
