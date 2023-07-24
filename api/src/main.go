package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	Id           int64  `db:"id"`
	Email        string `db:"email"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	DisplayName  string `db:"display_name"`
	Country      string `db:"country"`
	Geolocation  string `db:"geolocation"`
	PasswordHash string
	DateOfBirth  string `db:"date_of_birth"`
}

type UserId struct {
	Id int64 `db:"new_user_id"`
}

func main() {

	connStr := "user=binder_usr password=binder_best_app dbname=binder_all sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()

	router.GET("/app-revision", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"commitSha": "local"})
	})

	router.GET("/user/list", func(c *gin.Context) {
		users := []User{}
		err = db.Select(&users, "SELECT id, email, first_name, last_name, display_name, country, geolocation, date_of_birth FROM public.users")

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, users)
	})

	router.POST("/user", func(c *gin.Context) {
		user := User{}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		row := db.QueryRowx("SELECT sp_create_user($1, $2, $3, $4, $5, $6, $7, $8, $9)", user.Email, user.PasswordHash, user.FirstName, user.LastName, user.DisplayName, user.DateOfBirth, user.Country, -41, -89)
		result := make(map[string]interface{})
		_ = row.MapScan(result)
		user.Id = result["sp_create_user"].(int64)
		c.JSON(http.StatusCreated, user)
	})

	router.Run(":8080")
}
