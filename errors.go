package dari

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")
var ErrNotSupportNoSk = errors.New("not supported, no sk on table")
