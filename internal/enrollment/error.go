package enrollment

import (
	"errors"
	"fmt"
)

var ErrUserIdRequired = errors.New("user id is required")
var ErrCourseIdRequired = errors.New("course id name is required")
var ErrStatusRequired = errors.New("status is required")

type ErrNotFound struct {
	EnrollmentsID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user '%s' doesn't exists", e.EnrollmentsID)
}

type ErrInvalidStatus struct {
	Status string
}

func (e ErrInvalidStatus) Error() string {
	return fmt.Sprintf("status '%s' is invalid", e.Status)
}
