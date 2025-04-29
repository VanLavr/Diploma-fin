package query

import (
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type GetStudentsFilters struct {
	Limit  int64
	Offset int64
	IDs    []string
	Emails []string
}

func (this GetStudentsFilters) Validate() error {
	if len(this.Emails) == 0 && len(this.IDs) == 0 {
		return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
	}
	for _, email := range this.Emails {
		if email == "" {
			return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
		}
	}
	for _, id := range this.IDs {
		if id == "" {
			return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
		}
	}

	return nil
}
