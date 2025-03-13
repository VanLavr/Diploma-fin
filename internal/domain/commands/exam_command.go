package commands

import (
	"time"

	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type UpdateDebtByID struct {
	DebtID      int64
	Date        time.Time
	TeacherUUID string
	StudentUUID string
}

func (this UpdateDebtByID) Validate() error {
	if this.DebtID == 0 {
		return log.ErrorWrapper(errors.ErrInvalidCommand, errors.ERR_DOMAIN, "")
	}

	return nil
}
