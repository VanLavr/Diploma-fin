package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"

	"github.com/VanLavr/Diploma-fin/internal/services/logic"
)

type FileHandler struct {
	fUsecase logic.FileUsecase
}

func NewFileHandler(fUsecase logic.FileUsecase) *FileHandler {
	return &FileHandler{fUsecase: fUsecase}
}

func (fh FileHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/file/upload", fh.ParsFile)
}

func (fh FileHandler) ParsFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from form: " + err.Error()})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file: " + err.Error()})
		return
	}
	defer f.Close()

	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read Excel file: " + err.Error()})
		return
	}
	defer excelFile.Close()

	if err := fh.fUsecase.ParseFile(c.Request.Context(), excelFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Excel data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": nil})
}
