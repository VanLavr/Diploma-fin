package errors

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Send500Or400(err error, c *gin.Context) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
