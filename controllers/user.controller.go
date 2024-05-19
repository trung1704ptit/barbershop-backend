package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"barbershop-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

	user := &models.User{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Photo:     currentUser.Photo,
		Position:  currentUser.Position,
		Intro:     currentUser.Intro,
		Phone:     currentUser.Phone,
		Birthday:  currentUser.Birthday,
		Roles:     currentUser.Roles,
		Provider:  currentUser.Provider,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": user}})
}

func (uc *UserController) GetUserByPhone(ctx *gin.Context) {
	phone := ctx.Param("phone")
	var user models.User

	result := uc.DB.Preload("Services").Preload("Points").Preload("ServicesHistory").First(&user, "phone = ?", phone)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Tài khoản không tồn tại"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
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
		Photo:     payload.Photo,
		Roles:     payload.Roles,
		Intro:     payload.Intro,
		Birthday:  payload.Birthday,
		UpdatedAt: now,
		CreatedAt: updatedUser.CreatedAt,
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

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedUser})
}

func (uc *UserController) FindUsers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "100")
	var roles = ctx.DefaultQuery("role", "")
	roleArray := pq.StringArray(strings.Split(roles, ","))

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
	results := uc.DB.Preload("Services").Preload("Points").Preload("ServicesHistory").Limit(intLimit).Offset(offset).Order("updated_at DESC")
	if len(roleArray) != 0 && roles != "" {
		results = results.Where("roles @> ?", roleArray)
	}
	results = results.Find(&users)
	// var userResults []models.UserResponse

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(users), "data": users})
}

func (bc *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	var user models.User
	if err := bc.DB.Preload("Services").First(&user, "id = ?", userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "User not found"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	// Start a transaction
	tx := bc.DB.Begin()

	// Remove the association with services
	if err := tx.Model(&user).Association("Services").Clear(); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not remove service associations"})
		return
	}

	// Delete the user
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not delete booking"})
		return
	}

	// Commit the transaction
	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deleted successfully"})
}

func (uc *UserController) GetAllUsers() ([]models.User, error) {
	var users []models.User
	results := uc.DB.Limit(10000).Find(&users)

	if results.Error != nil {
		return []models.User{}, results.Error
	}
	return users, nil
}
