package query

import (
	"github.com/VanLavr/Diploma-fin/utils/errors"
)

type GetDebtsFilters struct {
	StudentUUIDs []string
	TeacherUUIDs []string
	ExamIDs      []int64
	DebtIDs      []int64
	Limit        int64
	Offset       int64
}

func (this *GetDebtsFilters) Validate() error {
	for _, id := range this.StudentUUIDs {
		if id == "" {
			return errors.ErrInvalidFilters
		}
	}
	for _, id := range this.TeacherUUIDs {
		if id == "" {
			return errors.ErrInvalidFilters
		}
	}
	for _, id := range this.ExamIDs {
		if id == 0 {
			return errors.ErrInvalidFilters
		}
	}

	return nil
}

type GetExamsFilters struct {
	Limit  int64
	Offset int64
	IDs    []int64
}

func (this *GetExamsFilters) Validate() error {
	for _, id := range this.IDs {
		if id == 0 {
			return errors.ErrInvalidFilters
		}
	}

	return nil
}
