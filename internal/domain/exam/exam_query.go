package exam

import "github.com/VanLavr/Diploma-fin/pkg/errors"

type GetAllDebtsFilter struct {
	StudentUUID string
}

func (this *GetAllDebtsFilter) Validate() error {
	if this.StudentUUID == "" {
		return errors.ErrInvalidFilters
	}

	return nil
}
