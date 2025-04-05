package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/VanLavr/Diploma-fin/internal/controllers/dto"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/utils/auth"
)

type AuthHandler struct {
	teacherUsecase logic.TeacherUsecase
	studentUsecase logic.StudentUsecase
	auth           *auth.AuthMiddleware
}

func NewAuthHandler(teacherUsecase logic.TeacherUsecase, studentUsecase logic.StudentUsecase, auth *auth.AuthMiddleware) *AuthHandler {
	return &AuthHandler{
		teacherUsecase: teacherUsecase,
		studentUsecase: studentUsecase,
		auth:           auth,
	}
}

func (a AuthHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/teacher/login", a.TeacherLogin)
	group.GET("/student/login", a.StudentLogin)
}

func (a AuthHandler) TeacherLogin(c *gin.Context) {
	var request dto.LoginDTO
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	fmt.Println(request)

	teachers, err := a.teacherUsecase.GetTeacherByEmail(c.Request.Context(), request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if len(teachers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	token, err := a.auth.GenerateAccessToken(teachers[0].UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

func (a AuthHandler) StudentLogin(c *gin.Context) {
	var request dto.LoginDTO
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	fmt.Println(request)

	students, err := a.studentUsecase.GetStudentByEmail(c.Request.Context(), request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if len(students) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	token, err := a.auth.GenerateAccessToken(students[0].UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}
