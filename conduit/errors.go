package conduit

import "errors"

var ErrDuplicateEmail = errors.New("duplicate email")
var ErrDuplicateUsername = errors.New("duplicate username")
var ErrNotFound = errors.New("record not found")
var ErrUnAuthorized = errors.New("unauthorized")
