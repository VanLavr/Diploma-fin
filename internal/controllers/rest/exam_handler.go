package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/VanLavr/Diploma-fin/internal/controllers/dto"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/utils/auth"
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
	group.POST("/exam", this.CreateExam)                 // + admin
	group.PUT("/exam", this.UpdateExam)                  // + admin
	group.DELETE("/exam/:id", this.DeleteExam)           // + admin
	group.GET("/exam/all/:limit/:offset", this.GetExams) // + admin
	group.GET("/exam/:id", this.GetExam)                 // + admin

	group.DELETE("/debt/:id", this.DeleteDebt)           // + admin
	group.PUT("/debt", this.UpdateDebt)                  // + admin
	group.POST("/debt", this.CreateDebt)                 // + admin
	group.GET("/debt/:id", this.GetDebt)                 // + admin
	group.GET("/debt/all/:limit/:offset", this.GetDebts) // + admin
}

func (e ExamHandler) GetDebt(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	debt, err := e.examUsecase.GetDebt(c.Request.Context(), int64(id))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := dto.DebtDTOFromTypes(*debt)

	c.JSON(http.StatusOK, dto.GetDebtDTO{
		Err:  nil,
		Data: result,
	})
}

func (e ExamHandler) GetExam(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	exam, err := e.examUsecase.GetExam(c.Request.Context(), int64(id))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := dto.ExamDTOFromTypes(*exam)

	c.JSON(http.StatusOK, dto.GetExamDTO{
		Err:  nil,
		Data: result,
	})
}

func (this ExamHandler) GetDebts(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

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

	debts, err := this.examUsecase.GetDebts(c.Request.Context(), int64(limit), int64(offset))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := make([]dto.Debt, 0, len(debts))
	for _, debt := range debts {
		result = append(result, dto.DebtDTOFromTypes(debt))
	}

	c.JSON(http.StatusOK, dto.GetAllDebtsDTO{
		Err:  nil,
		Data: result,
	})
}

func (this ExamHandler) GetExams(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

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

	result := make([]dto.Exam, 0, len(exams))
	for _, exam := range exams {
		result = append(result, dto.ExamDTOFromTypes(exam))
	}

	c.JSON(http.StatusOK, dto.GetAllExamsDTO{
		Err:  nil,
		Data: result,
	})
}

func (this ExamHandler) DeleteExam(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

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

func (this ExamHandler) DeleteDebt(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := this.examUsecase.DeleteDebt(c.Request.Context(), int64(id)); err != nil {
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

func (this ExamHandler) UpdateDebt(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	var r dto.UpdateDebtDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	debtTypes, err := dto.TypesDebtFromUpdateDebtDTO(r)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := this.examUsecase.UpdateDebt(c.Request.Context(), *debtTypes); err != nil {
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

func (this ExamHandler) UpdateExam(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

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
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

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

func (this ExamHandler) CreateDebt(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	var r dto.CreateDebtDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	debt, err := dto.TypesDebtFromCreateDebtDTO(r)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	fmt.Println("DEBUG:", debt, debt.Exam.ID)

	id, err := this.examUsecase.CreateDebt(c.Request.Context(), *debt)
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
