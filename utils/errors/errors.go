package errors

import "errors"

type ErrorType string

const (
	ERR_DOMAIN         ErrorType = "error in DOMAIN layer"
	ERR_APPLICATION    ErrorType = "error in APPLICATION layer"
	ERR_INTERFACES     ErrorType = "error in INTERFACES layer"
	ERR_INFRASTRUCTURE ErrorType = "error in INFRASTRUCTURE layer"
)

var ErrInvalidFilters = errors.New("search filters are invalid")
var ErrInvalidCommand = errors.New("insert, delete or update command is invalid")
var ErrInvalidData = errors.New("provided data is invalid")
var ErroNoItemsFound = errors.New("no items found")
var ErroInvalidPassword = errors.New("wrong password")
var ErrUserDoesNotHaveRights = errors.New("user doesn't have required rights")
var ErrCannotDetermineUUID = errors.New("token is malformed: can't determine structure")

const MethodKey string = "in method"
