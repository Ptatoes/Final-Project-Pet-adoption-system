package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdoptionsRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/adoptions", func(c *gin.Context) { ListAdoptions(c, db) })
	r.POST("/adoptions/assign", func(c *gin.Context) { AssignAdoption(c, db) })
}

func ListAdoptions(c *gin.Context, db *sql.DB) {
	// Fetch adoption history
	adoptionsRows, err := db.Query(`SELECT a.id, p.name AS pet_name, ad.name AS adopter_name, a.adopted_at 
                                  FROM adoptions a
                                  JOIN pets p ON a.pet_id = p.id
                                  JOIN adopters ad ON a.adopter_id = ad.id`)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch adoption history")
		return
	}
	defer adoptionsRows.Close()

	var adoptions []map[string]interface{}
	for adoptionsRows.Next() {
		var id int
		var petName, adopterName, adoptedAt string
		adoptionsRows.Scan(&id, &petName, &adopterName, &adoptedAt)

		adoptions = append(adoptions, map[string]interface{}{
			"id":           id,
			"pet_name":     petName,
			"adopter_name": adopterName,
			"adopted_at":   adoptedAt,
		})
	}

	// Fetch available pets
	petsRows, err := db.Query("SELECT id, name, species FROM pets WHERE status != 'Adopted'")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch available pets")
		return
	}
	defer petsRows.Close()

	var pets []map[string]interface{}
	for petsRows.Next() {
		var id int
		var name, species string
		petsRows.Scan(&id, &name, &species)

		pets = append(pets, map[string]interface{}{
			"id":      id,
			"name":    name,
			"species": species,
		})
	}

	// Fetch adopters
	adoptersRows, err := db.Query("SELECT id, name FROM adopters")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch adopters")
		return
	}
	defer adoptersRows.Close()

	var adopters []map[string]interface{}
	for adoptersRows.Next() {
		var id int
		var name string
		adoptersRows.Scan(&id, &name)

		adopters = append(adopters, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	// Render the adoptions page
	c.HTML(http.StatusOK, "adoptions.html", gin.H{
		"adoptions": adoptions,
		"pets":      pets,
		"adopters":  adopters,
	})
}

func AssignAdoption(c *gin.Context, db *sql.DB) {
	// Assign adoption
	petID, _ := strconv.Atoi(c.PostForm("pet_id"))
	adopterID, _ := strconv.Atoi(c.PostForm("adopter_id"))
	_, err := db.Exec("INSERT INTO adoptions (pet_id, adopter_id) VALUES (?, ?)", petID, adopterID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to assign adoption")
		return
	}

	// Update pet status to adopted
	_, err = db.Exec("UPDATE pets SET status = 'Adopted' WHERE id = ?", petID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to update pet status")
		return
	}

	c.Redirect(http.StatusSeeOther, "/adoptions")
}
