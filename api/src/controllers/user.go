package controllers

import (
	"binder_api/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserController struct {
	db *sqlx.DB
}

type UserInterestRequest struct {
	UserId    int64
	Interests []string
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

func (controller UserController) RegisterUserEndpoints(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/user/list", controller.GetUserList)
	api.POST("/user", controller.CreateUser)
	api.PATCH("/user-interests", controller.UpdateUserInterests)
}

func (controller UserController) GetUserList(c *gin.Context) {
	users := []db.User{}
	err := controller.db.Select(&users, "SELECT * FROM users_info")

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, users)
}

func (controller UserController) CreateUser(c *gin.Context) {
	req := CreateUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	var userID int64
	user := db.User{}
	controller.db.Get(&userID, "SELECT sp_create_user($1, $2, $3, $4, $5, $6, $7, $8, $9)", req.Email, req.PasswordHash, req.FirstName, req.LastName, req.DisplayName, req.DateOfBirth, req.Country, -41, -89)
	controller.db.Get(&user, "SELECT * FROM users_info WHERE id = $1", userID)
	c.JSON(http.StatusCreated, user)
}

func (controller UserController) UpdateUserInterests(c *gin.Context) {
	req := UserInterestRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	var isOk bool
	controller.db.Get(&isOk, "SELECT sp_update_user_interests($1, $2)", req.UserId, pq.Array(req.Interests))
	user := db.User{}
	controller.db.Get(&user, "SELECT * FROM users_info WHERE id = $1", req.UserId)
	c.JSON(http.StatusOK, user)
}

func ProvideUserController(db *sqlx.DB) *UserController {
	return &UserController{db: db}
}
