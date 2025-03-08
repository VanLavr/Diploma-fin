package http

import (
	"net/http"

	"github.com/VanLavr/Diploma-fin/internal/application/student"
	"github.com/gin-gonic/gin"
)

type ChatroomHandler struct {
	usecase student.StudentUsecase
	gin     *gin.RouterGroup
}

func NewChatroomHandler() *ChatroomHandler {
	return &ChatroomHandler{}
}

func (this ChatroomHandler) RegisterRoutes() {
	this.gin.GET("/all_debts/:UUID", this.getAllDebts)
	this.gin.POST("/notification/:UUID/:examID", this.sendNotification)
}

func (this ChatroomHandler) sendNotification(c *gin.Context) {
	panic("not implemented")
}

func (this ChatroomHandler) getAllDebts(c *gin.Context) {
	debts, err := this.usecase.GetAllDebts(c.Request.Context(), c.Param("UUID"))
	switch {
	case err == nil:
	default:
		c.JSON(http.StatusInternalServerError, getAllDebtsDTO{
			Err:  err,
			Data: nil,
		})
		return
	}

	exams := make([]Exam, len(debts))
	for i, exam := range debts {
		copyExam(&exams[i], &exam)
	}

	c.JSON(http.StatusOK, getAllDebtsDTO{
		Err:  err,
		Data: exams,
	})
}
