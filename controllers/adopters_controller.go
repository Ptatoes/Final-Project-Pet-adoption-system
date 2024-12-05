package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdoptersRoutes defines the routes for managing adopters
func AdoptersRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/adopters", func(c *gin.Context) { ListAdopters(c, db) })
	r.POST("/adopters/add", func(c *gin.Context) { AddAdopter(c, db) })
	r.POST("/adopters/delete", func(c *gin.Context) { DeleteAdopter(c, db) })
	r.GET("/adopters/edit/:id", func(c *gin.Context) { ShowEditForm(c, db) }) // Route to show edit form
	r.POST("/adopters/edit/:id", func(c *gin.Context) { EditAdopter(c, db) }) // Route to save edited adopter
}

// ListAdopters fetches and displays a list of adopters
func ListAdopters(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT id, name, contact FROM adopters")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch adopters"})
		return
	}
	defer rows.Close()

	var adopters []map[string]interface{}
	for rows.Next() {
		var id int
		var name, contact string
		if err := rows.Scan(&id, &name, &contact); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan adopters"})
			return
		}
		adopters = append(adopters, map[string]interface{}{
			"id":      id,
			"name":    name,
			"contact": contact,
		})
	}

	// Pass success message if set (for display in HTML)
	successMessage := c.DefaultQuery("success", "")

	c.HTML(http.StatusOK, "adopters.html", gin.H{
		"adopters":       adopters,
		"successMessage": successMessage,
	})
}

// AddAdopter adds a new adopter to the database
func AddAdopter(c *gin.Context, db *sql.DB) {
	name := c.PostForm("name")
	contact := c.PostForm("contact")
	_, err := db.Exec("INSERT INTO adopters (name, contact) VALUES (?, ?)", name, contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add adopter"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/adopters?success=Adopter added successfully!")
}

// DeleteAdopter deletes an adopter from the database
func DeleteAdopter(c *gin.Context, db *sql.DB) {
	id := c.PostForm("id")
	_, err := db.Exec("DELETE FROM adopters WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete adopter"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/adopters?success=Adopter deleted successfully!")
}

// ShowEditForm displays the form for editing an adopter's details
func ShowEditForm(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	row := db.QueryRow("SELECT id, name, contact FROM adopters WHERE id = ?", id)

	var adopter struct {
		ID      int
		Name    string
		Contact string
	}
	if err := row.Scan(&adopter.ID, &adopter.Name, &adopter.Contact); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Adopter not found"})
		return
	}

	c.HTML(http.StatusOK, "edit_adopter.html", gin.H{"adopter": adopter})
}

// EditAdopter updates the adopter's details in the database
func EditAdopter(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	name := c.PostForm("name")
	contact := c.PostForm("contact")

	_, err := db.Exec("UPDATE adopters SET name = ?, contact = ? WHERE id = ?", name, contact, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update adopter"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/adopters?success=Adopter updated successfully!")
}
