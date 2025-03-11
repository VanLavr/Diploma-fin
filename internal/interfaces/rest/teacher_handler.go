package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/VanLavr/Diploma-fin/internal/application/logic"
	"github.com/VanLavr/Diploma-fin/internal/interfaces/dto"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
)

type TeacherHandler struct {
	teacherUsecase logic.TeacherUsecase
	studentUsecase logic.StudentUsecase
	gin            *gin.RouterGroup
}

func NewTeacherHandler() *TeacherHandler {
	return &TeacherHandler{}
}

func (this TeacherHandler) RegisterRoutes() {
	this.gin.GET("/all_debts/:UUID", this.getAllDebts)
	this.gin.POST("/set_date", this.setDate)
}

func (this TeacherHandler) setDate(c *gin.Context) {
	var request dto.SetDebtDateDTO
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	switch err := this.teacherUsecase.SetDate(c.Request.Context(), request.TeacherUUID, request.ExamDate, request.DebtID); err {
	case nil:
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
}

func (this TeacherHandler) getAllDebts(c *gin.Context) {
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
