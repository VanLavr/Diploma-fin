package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/VanLavr/Diploma-fin/internal/controllers/dto"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type ExamHandler struct {
	examUsecase logic.ExamUsecase
}

func NewExamHandler(examUsecase logic.ExamUsecase) *ExamHandler {
	return &ExamHandler{
		examUsecase: examUsecase,
	}
}

func (this ExamHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/exam", this.CreateExam)
	group.PUT("/exam", this.UpdateExam)
	group.DELETE("/exam/:id", this.DeleteExam)
	group.GET("exam/all/:limit/:offset", this.GetExams)
}

func (this ExamHandler) GetExams(c *gin.Context) {
	lim := c.Param("limit")
	limit, err := strconv.Atoi(lim)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	offsetStr := c.Param("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	exams, err := this.examUsecase.GetExams(c.Request.Context(), int64(limit), int64(offset))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := make([]dto.Debt, 0, len(exams))
	for _, exam := range exams {
		result = append(result, dto.DebtDTOFromTypes(exam))
	}

	c.JSON(http.StatusOK, dto.GetAllDebtsDTO{
		Err:  nil,
		Data: result,
	})
}

func (this ExamHandler) DeleteExam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := this.examUsecase.DeleteExam(c.Request.Context(), int64(id)); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetAllDebtsDTO{
		Err:  nil,
		Data: nil,
	})
}

func (this ExamHandler) UpdateExam(c *gin.Context) {
	var r dto.UpdateExamDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := this.examUsecase.UpdateExam(c.Request.Context(), dto.TypesExamFromUpdateExamDTO(r)); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetAllDebtsDTO{
		Err:  nil,
		Data: nil,
	})
}

func (this ExamHandler) CreateExam(c *gin.Context) {
	var r dto.CreateExamDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	id, err := this.examUsecase.CreateExam(c.Request.Context(), dto.TypesExamFromCreateExamDTO(r))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, dto.CreateExamResponseDTO{
		Err:  nil,
		Data: int(id),
	})
}
