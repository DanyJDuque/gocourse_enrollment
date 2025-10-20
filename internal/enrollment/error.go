package enrollment

import (
	"errors"
	"fmt"
)

var ErrUserIdRequiered = errors.New("user id is requiered")
var ErrCourseIdRequiered = errors.New("course id name is requiered")

type ErrNotFound struct {
	UserID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user '%s' doesn't exist", e.UserID)
}
