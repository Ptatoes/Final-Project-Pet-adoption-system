package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdoptersRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/adopters", func(c *gin.Context) { ListAdopters(c, db) })
	r.POST("/adopters/add", func(c *gin.Context) { AddAdopter(c, db) })
	r.POST("/adopters/delete", func(c *gin.Context) { DeleteAdopter(c, db) })
}

func ListAdopters(c *gin.Context, db *sql.DB) {
	rows, _ := db.Query("SELECT id, name, contact FROM adopters")
	defer rows.Close()

	var adopters []map[string]interface{}
	for rows.Next() {
		var id int
		var name, contact string
		rows.Scan(&id, &name, &contact)

		adopters = append(adopters, map[string]interface{}{
			"id":      id,
			"name":    name,
			"contact": contact,
		})
	}
	c.HTML(http.StatusOK, "adopters.html", gin.H{"adopters": adopters})
}

func AddAdopter(c *gin.Context, db *sql.DB) {
	name := c.PostForm("name")
	contact := c.PostForm("contact")
	_, _ = db.Exec("INSERT INTO adopters (name, contact) VALUES (?, ?)", name, contact)
	c.Redirect(http.StatusSeeOther, "/adopters")
}

func DeleteAdopter(c *gin.Context, db *sql.DB) {
	id := c.PostForm("id")
	_, _ = db.Exec("DELETE FROM adopters WHERE id = ?", id)
	c.Redirect(http.StatusSeeOther, "/adopters")
}
