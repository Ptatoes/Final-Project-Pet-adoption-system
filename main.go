package main

import (
	"database/sql"
	"pet_adoption/controllers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pet_adoption")
	if err != nil {
		panic(err)
	}
}

func main() {
	initDB()

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", controllers.HomePage)
	r.GET("/FAQ", controllers.FAQPage)

	controllers.PetsRoutes(r, db)
	controllers.AdoptersRoutes(r, db)
	controllers.AdoptionsRoutes(r, db)

	r.Run(":8080")
}
