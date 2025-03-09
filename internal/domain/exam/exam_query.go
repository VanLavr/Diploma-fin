package exam

import "github.com/VanLavr/Diploma-fin/pkg/errors"

type GetDebtsFilters struct {
	StudentUUIDs []string
	TeacherUUIDs []string
	ExamIDs      []int64
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

	if len(this.StudentUUIDs) == 0 && len(this.TeacherUUIDs) == 0 && len(this.ExamIDs) == 0 {
		return errors.ErrInvalidFilters
	}

	return nil
}

type GetExamsFilters struct {
	IDs []int64
}

func (this *GetExamsFilters) Validate() error {
	for _, id := range this.IDs {
		if id == 0 {
			return errors.ErrInvalidFilters
		}
	}

	return nil
}
