package errors

import "errors"

var ErrInvalidFilters = errors.New("search filters are invalid")
var ErroNoItemsFound = errors.New("no items found")
