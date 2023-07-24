package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	Id          int64  `db:"id"`
	Email       string `db:"email"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	DisplayName string `db:"display_name"`
	Country     string `db:"country"`
	Geolocation string `db:"geolocation"`
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
		err = db.Select(&users, "SELECT id, email, first_name, last_name, display_name, country, geolocation FROM public.users")

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, users)
	})

	router.Run(":8080")
}
