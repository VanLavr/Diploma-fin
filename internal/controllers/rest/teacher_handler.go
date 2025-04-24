package rest

import (
	"fmt"
	"net/http"
	"strconv"

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
	group.GET("/teacher/all_debts/:UUID", this.getAllDebts)
	group.POST("/set_date", this.setDate)
	group.POST("/teacher", this.CreateTeacher)
	group.PUT("/teacher", this.UpdateTeacher)
	group.DELETE("/teacher/:uuid", this.DeleteTeacher)
	group.GET("/teacher/all/:limit/:offset", this.GetTeachers)
}

func (t TeacherHandler) CreateTeacher(c *gin.Context) {
	var r dto.CreateTeacherDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	UUID, err := t.teacherUsecase.CreateTeacher(c.Request.Context(), dto.TypesTeacherFromCreateTeacherDTO(r))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, dto.CreateTeacherResponseDTO{
		Err:  nil,
		Data: UUID,
	})
}

func (t TeacherHandler) UpdateTeacher(c *gin.Context) {
	var r dto.UpdateTeacherDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := t.teacherUsecase.UpdateTeacher(c.Request.Context(), dto.TypesTeacherFromUpdateTeachertDTO(r)); err != nil {
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

func (t TeacherHandler) DeleteTeacher(c *gin.Context) {
	if err := t.teacherUsecase.DeleteTeacher(c.Request.Context(), c.Param("uuid")); err != nil {
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

func (t TeacherHandler) GetTeachers(c *gin.Context) {
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

	teachers, err := t.teacherUsecase.GetTeachers(c.Request.Context(), int64(limit), int64(offset))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := make([]dto.Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		result = append(result, dto.TeacherDTOFromTypes(teacher))
	}

	c.JSON(http.StatusOK, dto.GetAllTeachersDTO{
		Err:  nil,
		Data: result,
	})
}

func (this TeacherHandler) setDate(c *gin.Context) {
	var request dto.SetDebtDateDTO
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	fmt.Println(request)

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
	debts, err := this.teacherUsecase.GetAllDebts(c.Request.Context(), c.Param("UUID"))
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
