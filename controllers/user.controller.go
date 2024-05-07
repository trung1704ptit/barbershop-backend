package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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

	result := uc.DB.Preload("Services").Preload("Points").Preload("ServicesHistory").First(&user, "phone = ?", phone)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Tài khoản không tồn tại"})
		return
	}

	userResponse := &models.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Photo:           user.Photo,
		Phone:           user.Phone,
		Birthday:        user.Birthday,
		Points:          user.Points,
		Services:        user.Services,
		ServicesHistory: user.ServicesHistory,
		Role:            user.Role,
		Provider:        user.Provider,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": userResponse})
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	var payload *models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedUser models.User
	result := uc.DB.First(&updatedUser, "id = ?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Tài khoản user không tồn tại"})
		return
	}

	now := time.Now()
	userToUpdate := models.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Phone:     payload.Phone,
		Birthday:  payload.Birthday,
		UpdatedAt: now,
	}

	userResponse := &models.UserResponse{
		ID:              updatedUser.ID,
		Name:            updatedUser.Name,
		Email:           updatedUser.Email,
		Photo:           updatedUser.Photo,
		Phone:           updatedUser.Phone,
		Birthday:        updatedUser.Birthday,
		ServicesHistory: updatedUser.ServicesHistory,
		Role:            updatedUser.Role,
		Provider:        updatedUser.Provider,
		CreatedAt:       updatedUser.CreatedAt,
		UpdatedAt:       updatedUser.UpdatedAt,
	}

	updateResult := uc.DB.Model(&updatedUser).Updates(userToUpdate)
	if updateResult.Error != nil {
		errorMsg := updateResult.Error.Error()
		if strings.Contains(errorMsg, "duplicate key value") && strings.Contains(errorMsg, "idx_users_email") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Địa chỉ email đã tồn tại."})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": userResponse})
}

func (uc *UserController) FindUsers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "100")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid page parameter"})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid limit parameter"})
		return
	}

	offset := (intPage - 1) * intLimit

	var users []models.User
	results := uc.DB.Preload("Services").Preload("Points").Preload("ServicesHistory").Limit(intLimit).Offset(offset).Find(&users)
	var userResults []models.UserResponse

	for _, user := range users {
		userResults = append(userResults, models.UserResponse{
			ID:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			Photo:           user.Photo,
			Phone:           user.Phone,
			Birthday:        user.Birthday,
			Services:        user.Services,
			Points:          user.Points,
			ServicesHistory: user.ServicesHistory,
			Role:            user.Role,
			Provider:        user.Provider,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		})
	}

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(userResults), "data": userResults})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	result := uc.DB.Delete(&models.User{}, "id = ?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No User with that id exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (uc *UserController) GetAllUsers() ([]models.UserResponse, error) {
	var users []models.User
	results := uc.DB.Limit(10000).Find(&users)

	if results.Error != nil {
		return []models.UserResponse{}, results.Error
	}

	var userResults []models.UserResponse

	for _, user := range users {
		userResults = append(userResults, models.UserResponse{
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
		})
	}
	return userResults, nil
}
