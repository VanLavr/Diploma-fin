package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/VanLavr/Diploma-fin/internal/controllers/dto"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type TeacherHandler struct {
	teacherUsecase logic.TeacherUsecase
	studentUsecase logic.StudentUsecase
}

func NewTeacherHandler(teacherUsecase logic.TeacherUsecase, studentUsecase logic.StudentUsecase) *TeacherHandler {
	return &TeacherHandler{
		teacherUsecase: teacherUsecase,
		studentUsecase: studentUsecase,
	}
}

func (this TeacherHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/all_debts/:UUID", this.getAllDebts)
	group.POST("/set_date", this.setDate)
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
