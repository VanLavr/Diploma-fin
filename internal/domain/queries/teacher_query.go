package query

import (
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type GetTeachersFilters struct {
	UUIDs  []string
	Emails []string
	Limit  int64
	Offset int64
}

func (this GetTeachersFilters) Validate() error {
	if len(this.Emails) == 0 && len(this.UUIDs) == 0 && this.Limit == 0 {
		return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
	}
	for _, email := range this.Emails {
		if email == "" {
			return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
		}
	}
	for _, id := range this.UUIDs {
		if id == "" {
			return log.ErrorWrapper(errors.ErrInvalidFilters, errors.ERR_DOMAIN, "")
		}
	}
	return nil
}
