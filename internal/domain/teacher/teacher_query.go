package teacher

import "github.com/VanLavr/Diploma-fin/pkg/errors"

type GetTeachersFilters struct {
	UUIDs []string
}

func (this GetTeachersFilters) Validate() error {
	for _, id := range this.UUIDs {
		if id == "" {
			return errors.ErrInvalidFilters
		}
	}
	return nil
}
