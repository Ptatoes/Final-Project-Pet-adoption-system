package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

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
	r.GET("/", homePage)
	r.GET("/pets", listPets)
	r.POST("/pets/add", addPet)
	r.POST("/pets/update", updatePet)
	r.POST("/pets/delete", deletePet)
	r.GET("/adopters", listAdopters)
	r.POST("/adopters/add", addAdopter)
	r.POST("/adopters/delete", deleteAdopter)
	r.GET("/adoptions", listAdoptions)
	r.POST("/adoptions/assign", assignAdoption)
	r.GET("/FAQ", faqPage)
	r.GET("/pets/search", searchPets)

	r.Run(":8080")
}

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func faqPage(c *gin.Context) {
	c.HTML(http.StatusOK, "FAQ.html", nil)

}

func listPets(c *gin.Context) {
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
	c.HTML(http.StatusOK, "pets.html", gin.H{"pets": pets})
}

func addPet(c *gin.Context) {
	name := c.PostForm("name")
	species := c.PostForm("species")
	_, _ = db.Exec("INSERT INTO pets (name, species) VALUES (?, ?)", name, species)
	c.Redirect(http.StatusSeeOther, "/pets")
}

func updatePet(c *gin.Context) {
	id := c.PostForm("id")
	status := c.PostForm("status")
	_, _ = db.Exec("UPDATE pets SET status = ? WHERE id = ?", status, id)
	c.Redirect(http.StatusSeeOther, "/pets")
}

func deletePet(c *gin.Context) {
	id := c.PostForm("id")
	_, _ = db.Exec("DELETE FROM pets WHERE id = ?", id)
	c.Redirect(http.StatusSeeOther, "/pets")
}

func searchPets(c *gin.Context) {
	query := c.Query("query")
	log.Println("Search query:", query)

	rows, err := db.Query("SELECT id, name, species, status FROM pets WHERE name LIKE ? OR species LIKE ?", "%"+query+"%", "%"+query+"%")
	if err != nil {
		log.Println("Error executing query:", err)
		c.String(http.StatusInternalServerError, "Database query error")
		return
	}
	defer rows.Close()

	var pets []map[string]interface{}
	for rows.Next() {
		var id int
		var name, species, status string
		if err := rows.Scan(&id, &name, &species, &status); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		log.Printf("Fetched pet: ID=%d, Name=%s, Species=%s, Status=%s\n", id, name, species, status)
		pets = append(pets, map[string]interface{}{
			"id":      id,
			"name":    name,
			"species": species,
			"status":  status,
		})
	}
	log.Println("Fetched pets:", pets)

	c.HTML(http.StatusOK, "pets.html", gin.H{"pets": pets})
}

func listAdopters(c *gin.Context) {
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

func addAdopter(c *gin.Context) {
	name := c.PostForm("name")
	contact := c.PostForm("contact")
	_, _ = db.Exec("INSERT INTO adopters (name, contact) VALUES (?, ?)", name, contact)
	c.Redirect(http.StatusSeeOther, "/adopters")
}

func deleteAdopter(c *gin.Context) {
	id := c.PostForm("id")
	_, _ = db.Exec("DELETE FROM adopters WHERE id = ?", id)
	c.Redirect(http.StatusSeeOther, "/adopters")
}

func listAdoptions(c *gin.Context) {
	adoptionsRows, _ := db.Query(`SELECT a.id, p.name AS pet_name, ad.name AS adopter_name, a.adopted_at 
                                  FROM adoptions a
                                  JOIN pets p ON a.pet_id = p.id
                                  JOIN adopters ad ON a.adopter_id = ad.id`)
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

	petsRows, _ := db.Query("SELECT id, name, species FROM pets WHERE status != 'Adopted'")
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

	adoptersRows, _ := db.Query("SELECT id, name FROM adopters")
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

	c.HTML(http.StatusOK, "adoptions.html", gin.H{
		"adoptions": adoptions,
		"pets":      pets,
		"adopters":  adopters,
	})
}

func assignAdoption(c *gin.Context) {
	petID, _ := strconv.Atoi(c.PostForm("pet_id"))
	adopterID, _ := strconv.Atoi(c.PostForm("adopter_id"))
	_, _ = db.Exec("INSERT INTO adoptions (pet_id, adopter_id) VALUES (?, ?)", petID, adopterID)
	_, _ = db.Exec("UPDATE pets SET status = 'Adopted' WHERE id = ?", petID)
	c.Redirect(http.StatusSeeOther, "/adoptions")
}
