package student

import (
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
)

type GetStudentsFilters struct {
	IDs []string
}

func (this GetStudentsFilters) Validate() error {
	if len(this.IDs) == 0 {
		return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
	}
	for _, id := range this.IDs {
		if id == "" {
			return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
		}
	}

	return nil
}
