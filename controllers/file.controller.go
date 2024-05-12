package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileController struct {
	DB *gorm.DB
}

type FileData struct {
	FilePath string `json:"filePath"`
}

type FilesData struct {
	FilePaths string `json:"filePaths"`
}

func NewFileController(DB *gorm.DB) FileController {
	return FileController{DB}
}

func (fc *FileController) Upload(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": fmt.Sprintf("Error retrieving the file: %s", err.Error())})
		return
	}
	defer file.Close()

	// Ensure the "uploads" directory exists; create it if not
	uploadsDir := "uploads/"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadsDir, 0755)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Error creating uploads directory: %s", err.Error())})
			return
		}
	}

	// Create a new file in the uploads directory
	filePath := uploadsDir + fileHeader.Filename
	out, err := os.Create("uploads/" + fileHeader.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": fmt.Sprintf("Error creating the file: %s", err.Error())})
		return
	}
	defer out.Close()

	// Copy the file content from the uploaded file to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": fmt.Sprintf("Error saving the file: %s", err.Error())})
		return
	}

	data := FileData{
		FilePath: filePath,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (fc *FileController) MultiUpload(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Error retrieving form data: %s", err.Error())})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "No files uploaded"})
		return
	}

	// Ensure the "uploads" directory exists; create it if not
	uploadsDir := "uploads/"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadsDir, 0755)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Error creating uploads directory: %s", err.Error())})
			return
		}
	}

	var filepaths []string
	for _, file := range files {
		filePath := uploadsDir + file.Filename
		out, err := os.Create(filePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Error creating the file: %s", err.Error())})
			return
		}
		defer out.Close()

		uploadedFile, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Error opening the uploaded file: %s", err.Error())})
			return
		}
		defer uploadedFile.Close()

		_, err = io.Copy(out, uploadedFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Error saving the file: %s", err.Error())})
			return
		}

		filepaths = append(filepaths, filePath)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Files uploaded successfully", "filepaths": filepaths})
}
