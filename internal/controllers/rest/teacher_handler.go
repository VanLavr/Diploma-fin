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
	"github.com/VanLavr/Diploma-fin/utils/hasher"
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
	group.GET("/teacher/all_debts/:UUID", this.getAllDebts)    // + teacher
	group.POST("/set_date", this.setDate)                      // + teacher
	group.POST("/teacher", this.CreateTeacher)                 // + admin
	group.PUT("/teacher", this.UpdateTeacher)                  // + admin,teacher
	group.DELETE("/teacher/:uuid", this.DeleteTeacher)         // + admin
	group.GET("/teacher/all/:limit/:offset", this.GetTeachers) // + admin
	group.GET("/teacher/:uuid", this.GetTeacher)               // + admin,teacher
	group.PUT("/teacher/pass", this.UpdateTeacherPassword)     // + teacher
}

func (t TeacherHandler) UpdateTeacherPassword(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.TeacherRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

	var r dto.UpdateTeacherPasswordDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	teacher, err := t.teacherUsecase.GetTeacher(c.Request.Context(), r.UUID)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !hasher.Hshr.Validate(teacher.Password, r.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErroInvalidPassword.Error(),
		})
		return
	}

	if err := t.teacherUsecase.ChangePassword(c.Request.Context(), r.UUID, r.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"password": "changed"})
}

func (t TeacherHandler) GetTeacher(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole || c.Value(auth.RoleKey) != auth.TeacherRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

	uuid := c.Param("uuid")

	teacher, err := t.teacherUsecase.GetTeacher(c.Request.Context(), uuid)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := dto.TeacherDTOFromTypes(teacher)

	c.JSON(http.StatusOK, dto.GetTeacherDTO{
		Err:  nil,
		Data: result,
	})
}

func (t TeacherHandler) CreateTeacher(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

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
	if c.Value(auth.RoleKey) != auth.AdminRole || c.Value(auth.RoleKey) != auth.TeacherRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

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
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

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
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
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
	if c.Value(auth.RoleKey) != auth.TeacherRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

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
	if c.Value(auth.RoleKey) != auth.TeacherRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

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
