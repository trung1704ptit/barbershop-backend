package controllers

import (
	"net/http"

	"barbershop-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Photo:     currentUser.Photo,
		Phone:     currentUser.Phone,
		Birthday:  currentUser.Birthday,
		Role:      currentUser.Role,
		Provider:  currentUser.Provider,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (uc *UserController) GetUserByPhone(ctx *gin.Context) {
	phone := ctx.Param("phone")
	var user models.User
	result := uc.DB.First(&user, "phone = ?", phone)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "User not found"})
		return
	}
	userResponse := &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Photo:     user.Photo,
		Phone:     user.Phone,
		Birthday:  user.Birthday,
		Role:      user.Role,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": userResponse})
}
