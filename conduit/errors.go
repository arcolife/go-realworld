package conduit

import "errors"

var ErrDuplicateEmail = errors.New("duplicate email")
var ErrDuplicateUsername = errors.New("duplicate username")
var ErrUnAuthorized = errors.New("unauthorized")
