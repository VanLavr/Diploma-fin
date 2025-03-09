package student

import "github.com/VanLavr/Diploma-fin/pkg/errors"

type GetStudentsFilters struct {
	IDs []string
}

func (this GetStudentsFilters) Validate() error {
	for _, id := range this.IDs {
		if id == "" {
			return errors.ErrInvalidFilters
		}
	}

	return nil
}
