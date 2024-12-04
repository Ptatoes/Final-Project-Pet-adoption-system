package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FAQPage(c *gin.Context) {
	c.HTML(http.StatusOK, "FAQ.html", nil)
}

func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
