package controllers

import (
	"binder_api/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type UserController struct {
	db     *sqlx.DB
	logger *zap.Logger
}

type CreateUserRequest struct {
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	DisplayName  string
	DateOfBirth  string
	Country      string
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
	users := []db.User{}
	err := controller.db.Select(&users, "SELECT * FROM users_info")

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
	}
	user := db.User{}
	err := controller.db.Get(&user, "SELECT * FROM sp_create_user($1, $2, $3, $4, $5, $6, $7, $8, $9)", req.Email, req.PasswordHash, req.FirstName, req.LastName, req.DisplayName, req.DateOfBirth, req.Country, -41, -89)
	if err != nil {
		controller.logger.Error("CreateUser() failed",
			zap.Error(err),
			zap.String("user_email", req.Email))
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (controller UserController) UpdateUserInterests(c *gin.Context) {
	req := SetUserInterestRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	user := db.User{}
	if err := controller.db.Get(&user, "SELECT * FROM sp_update_user_interests($1, $2)",
		req.UserId, pq.Array(req.Interests)); err != nil {
		controller.logger.Error("UpdateUserInterests() failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{})
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

	user := db.User{}
	if err := controller.db.Get(&user, "SELECT * FROM sp_update_user_photos($1, $2)",
		req.UserId, pq.Array(req.PhotoUrls)); err != nil {
		controller.logger.Error("UpdateUserPhoto() failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{})
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

	var user db.User
	if err := controller.db.Get(&user, "SELECT * FROM sp_update_user_filters($1, $2, $3, $4, $5)",
		req.UserId, req.MinDistanceKm, req.MaxDistanceKm, req.MinAge, req.MaxAge); err != nil {
		controller.logger.Error("UpdateUserFilter() failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

func ProvideUserController(logger *zap.Logger, db *sqlx.DB) *UserController {
	return &UserController{db: db, logger: logger}
}
