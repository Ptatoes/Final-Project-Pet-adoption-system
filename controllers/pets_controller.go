package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PetsRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/pets", func(c *gin.Context) { ListPets(c, db) })
	r.POST("/pets/add", func(c *gin.Context) { AddPet(c, db) })
	r.POST("/pets/update", func(c *gin.Context) { UpdatePet(c, db) })
	r.POST("/pets/delete", func(c *gin.Context) { DeletePet(c, db) })
	r.GET("/pets/search", func(c *gin.Context) { SearchPets(c, db) })
}

func ListPets(c *gin.Context, db *sql.DB) {
	rows, _ := db.Query("SELECT id, name, species, status FROM pets")
	defer rows.Close()

	var pets []map[string]interface{}
	for rows.Next() {
		var id int
		var name, species, status string
		rows.Scan(&id, &name, &species, &status)

		pets = append(pets, map[string]interface{}{
			"id":      id,
			"name":    name,
			"species": species,
			"status":  status,
		})
	}

	// Get the success message, if any
	successMessage := c.Query("successMessage")

	c.HTML(http.StatusOK, "pets.html", gin.H{
		"pets":           pets,
		"successMessage": successMessage,
	})
}

func AddPet(c *gin.Context, db *sql.DB) {
	name := c.PostForm("name")
	species := c.PostForm("species")
	_, err := db.Exec("INSERT INTO pets (name, species) VALUES (?, ?)", name, species)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to add pet")
		return
	}
	c.Redirect(http.StatusSeeOther, "/pets?successMessage=Pet successfully added!")
}

func UpdatePet(c *gin.Context, db *sql.DB) {
	id := c.PostForm("id")
	status := c.PostForm("status")
	_, _ = db.Exec("UPDATE pets SET status = ? WHERE id = ?", status, id)
	c.Redirect(http.StatusSeeOther, "/pets?successMessage=Pet successfully Updated!")
}

func DeletePet(c *gin.Context, db *sql.DB) {
	id := c.PostForm("id")
	_, err := db.Exec("DELETE FROM pets WHERE id = ?", id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete pet")
		return
	}
	c.Redirect(http.StatusSeeOther, "/pets?successMessage=Pet successfully deleted!")
}

func SearchPets(c *gin.Context, db *sql.DB) {
	query := c.Query("query")
	rows, _ := db.Query("SELECT id, name, species, status FROM pets WHERE name LIKE ? OR species LIKE ?", "%"+query+"%", "%"+query+"%")
	defer rows.Close()

	var pets []map[string]interface{}
	for rows.Next() {
		var id int
		var name, species, status string
		rows.Scan(&id, &name, &species, &status)
		pets = append(pets, map[string]interface{}{
			"id":      id,
			"name":    name,
			"species": species,
			"status":  status,
		})
	}
	c.HTML(http.StatusOK, "pets.html", gin.H{"pets": pets})
}
