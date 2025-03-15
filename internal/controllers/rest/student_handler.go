package rest

import (
	"net/http"
	"strconv"

	"github.com/VanLavr/Diploma-fin/internal/controllers/dto"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentUsecase logic.StudentUsecase
}

func NewStudentHandler(studentUsecase logic.StudentUsecase) *StudentHandler {
	return &StudentHandler{
		studentUsecase: studentUsecase,
	}
}

func (this StudentHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/student/all_debts/:UUID", this.getAllDebts)
	group.POST("/notification/:UUID/:examID", this.sendNotification)
}

func (this StudentHandler) sendNotification(c *gin.Context) {
	id := c.Param("examID")
	examID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.ErrorWrapper(err, errors.ERR_INTERFACES, "wrong examen id")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	switch err := this.studentUsecase.SendNotification(c.Request.Context(), c.Param("UUID"), examID); err {
	case nil:
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
}

func (this StudentHandler) getAllDebts(c *gin.Context) {
	debts, err := this.studentUsecase.GetAllDebts(c.Request.Context(), c.Param("UUID"))
	switch {
	case err == nil:
	default:
		log.ErrorWrapper(err, errors.ERR_APPLICATION, "", "size of recieved data", len(debts))
		c.JSON(http.StatusInternalServerError, dto.GetAllDebtsDTO{
			Err:  err,
			Data: nil,
		})
		return
	}

	exams := make([]dto.Debt, len(debts))
	for i, exam := range debts {
		exams[i] = dto.DebtDTOFromTypes(exam)
	}

	c.JSON(http.StatusOK, dto.GetAllDebtsDTO{
		Err:  err,
		Data: exams,
	})
}
