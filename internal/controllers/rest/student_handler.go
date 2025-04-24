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
	group.POST("/student", this.CreateStudent)
	group.PUT("/student", this.UpdateStudent)
	group.DELETE("/student/:uuid", this.DeleteStduent)
	group.GET("student/all/:limit/:offset", this.GetStudents)
}

func (this StudentHandler) GetStudents(c *gin.Context) {
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
