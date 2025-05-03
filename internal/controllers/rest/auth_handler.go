package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/VanLavr/Diploma-fin/internal/controllers/dto"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/utils/auth"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/hasher"
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
	group.POST("/login", a.Login)

	permissions := group.Group("/perm/check", a.auth.ValidateAccessToken())
	permissions.GET("/current_permissions", a.GetCurrentPermissions)
}

func (a AuthHandler) GetCurrentPermissions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{auth.RoleKey: c.Value(auth.RoleKey)})
}

func (a AuthHandler) Login(c *gin.Context) {
	var request dto.LoginDTO
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	teachers, err := a.teacherUsecase.GetTeacherByEmail(c.Request.Context(), request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	students, err := a.studentUsecase.GetStudentByEmail(c.Request.Context(), request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if len(teachers) == 0 && len(students) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErroNoItemsFound,
		})
		return
	}

	if len(students) == 0 {
		fmt.Println(teachers[0].Password)
		if !hasher.Hshr.Validate(teachers[0].Password, request.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.ErroInvalidPassword,
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
	} else {
		if !hasher.Hshr.Validate(students[0].Password, request.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.ErroInvalidPassword,
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
}
