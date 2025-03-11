package query

import (
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
)

type GetTeachersFilters struct {
	UUIDs []string
}

func (this GetTeachersFilters) Validate() error {
	if len(this.UUIDs) == 0 {
		return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
	}
	for _, id := range this.UUIDs {
		if id == "" {
			return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
		}
	}
	return nil
}
