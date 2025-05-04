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

type StudentHandler struct {
	studentUsecase logic.StudentUsecase
}

func NewStudentHandler(studentUsecase logic.StudentUsecase) *StudentHandler {
	return &StudentHandler{
		studentUsecase: studentUsecase,
	}
}

func (this StudentHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/student/all_debts/:UUID", this.getAllDebts)          // + student
	group.POST("/notification/:UUID/:examID", this.sendNotification) // + student
	group.POST("/student", this.CreateStudent)                       // + admin
	group.PUT("/student", this.UpdateStudent)                        // + admin,student
	group.DELETE("/student/:uuid", this.DeleteStduent)               // + admin
	group.GET("/student/all/:limit/:offset", this.GetStudents)       // + admin
	group.GET("/student/:uuid", this.GetStudent)                     // + admin,student
	group.PUT("/student/pass", this.UpdatePassword)                  // + student
}

func (s StudentHandler) UpdatePassword(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.StudentRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

	var r dto.UpdateStudentPasswordDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	student, err := s.studentUsecase.GetStudent(c.Request.Context(), r.UUID)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !hasher.Hshr.Validate(student.Password, r.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErroInvalidPassword.Error(),
		})
		return
	}

	if err := s.studentUsecase.ChangePassword(c.Request.Context(), r.UUID, hasher.Hshr.Hash(r.NewPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"password": "changed"})
}

func (s StudentHandler) GetStudent(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole && c.Value(auth.RoleKey) != auth.StudentRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

	uuid := c.Param("uuid")

	student, err := s.studentUsecase.GetStudent(c.Request.Context(), uuid)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := dto.StudentDTOFromTypes(*student)

	c.JSON(http.StatusOK, dto.GetStudentDTO{
		Err:  nil,
		Data: result,
	})
}

func (this StudentHandler) GetStudents(c *gin.Context) {
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

	students, err := this.studentUsecase.GetStudents(c.Request.Context(), int64(limit), int64(offset))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := make([]dto.Student, 0, len(students))
	for _, student := range students {
		result = append(result, dto.StudentDTOFromTypes(student))
	}

	c.JSON(http.StatusOK, dto.GetAllStudentsDTO{
		Err:  nil,
		Data: result,
	})
}

func (this StudentHandler) DeleteStduent(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

	if err := this.studentUsecase.DeleteStudent(c.Request.Context(), c.Param("uuid")); err != nil {
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

func (this StudentHandler) UpdateStudent(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole && c.Value(auth.RoleKey) != auth.StudentRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

	var r dto.UpdateStudentDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := this.studentUsecase.UpdateStudent(c.Request.Context(), dto.TypeStudentFromUpdateStudentDTO(r)); err != nil {
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

func (this StudentHandler) CreateStudent(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

	var r dto.CreateStudentDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	UUID, err := this.studentUsecase.CreateStudent(c.Request.Context(), dto.TypesStudentFromCreateStudentDTO(r))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, dto.CreateStudentResponseDTO{
		Err:  nil,
		Data: UUID,
	})
}

func (this StudentHandler) sendNotification(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.StudentRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

	id := c.Param("examID")
	fmt.Println("1")
	examID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.ErrorWrapper(err, errors.ERR_INTERFACES, "wrong examen id")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	fmt.Println("2")

	switch err := this.studentUsecase.SendNotification(c.Request.Context(), c.Param("UUID"), examID); err {
	case nil:
	default:
		fmt.Println("3")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	fmt.Println("4")
	c.JSON(http.StatusOK, dto.GetAllDebtsDTO{
		Err:  nil,
		Data: nil,
	})
}

func (this StudentHandler) getAllDebts(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.StudentRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights.Error()})
		return
	}

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
