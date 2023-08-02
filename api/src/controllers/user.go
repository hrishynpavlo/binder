package controllers

import (
	"binder_api/db"
	"binder_api/workers"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	logger                *zap.Logger
	db                    *db.UserRepository
	userRegisteredChannel chan workers.UserRegisteredEvent
}

type CreateUserRequest struct {
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	DisplayName  string
	DateOfBirth  string
	Country      string
	Latitude     float64
	Longitude    float64
}

type SetUserInterestRequest struct {
	UserId    int64
	Interests []string
}

type SetUserPhotosRequest struct {
	UserId    int64
	PhotoUrls []string
}

type SetUserFiltersRequest struct {
	UserId        int64
	MinDistanceKm int8
	MaxDistanceKm int8
	MinAge        int8
	MaxAge        int8
}

func (controller UserController) RegisterUserEndpoints(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/user/list", controller.GetUserList)
	api.POST("/user", controller.CreateUser)
	api.PATCH("/user-interests", controller.UpdateUserInterests)
	api.PATCH("/user-photos", controller.UpdateUserPhoto)
	api.PATCH("/user-filters", controller.UpdateUserFilter)
}

func (controller UserController) GetUserList(c *gin.Context) {
	users, err := controller.db.GetAllUsers()

	if err != nil {
		controller.logger.Warn("GetUserList() failed, see error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "wrong request body"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (controller UserController) CreateUser(c *gin.Context) {
	req := CreateUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		controller.logger.Warn("CreateUser() wrong request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
		return
	}

	user, err := controller.db.CreateUser(req.Email, req.PasswordHash, req.FirstName, req.LastName, req.DisplayName, req.DateOfBirth, req.Country, req.Latitude, req.Longitude)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	event := workers.UserRegisteredEvent{UserId: user.Id, Latitude: req.Latitude, Longitude: req.Longitude}
	controller.userRegisteredChannel <- event
	c.JSON(http.StatusCreated, user)
}

func (controller UserController) UpdateUserInterests(c *gin.Context) {
	req := SetUserInterestRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	user, err := controller.db.UpdateUserInterests(req.UserId, req.Interests)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

func (controller UserController) UpdateUserPhoto(c *gin.Context) {
	req := SetUserPhotosRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := controller.db.UpdateUserPhoto(req.UserId, req.PhotoUrls)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

func (controller UserController) UpdateUserFilter(c *gin.Context) {
	var req SetUserFiltersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := controller.db.UpdateUserFilter(req.UserId, req.MinDistanceKm, req.MaxDistanceKm, req.MinAge, req.MaxAge)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

func ProvideUserController(logger *zap.Logger, db *db.UserRepository, userRegisteredChannel chan workers.UserRegisteredEvent) *UserController {
	return &UserController{logger: logger, db: db, userRegisteredChannel: userRegisteredChannel}
}
