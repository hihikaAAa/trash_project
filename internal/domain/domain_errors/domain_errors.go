// Package domainerrors
package domainerrors

import (
	"errors"
)

var(
	ErrTaskCanceled = errors.New("task was already canceled")
	ErrTaskDone = errors.New("task was already done")
)