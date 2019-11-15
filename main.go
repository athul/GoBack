package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Wubba Lubba Dub Dub",
		})
	})
	api := r.Group("/api")
	api.POST("/names/:nameID", nameAPI)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type name struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

var Names = []name{
	name{1, "Athul Cyriac Ajay"},
	name{2, "Elvis Jacob Ajay"},
	name{3, "Anton-Bivens-Davids"},
}

func nameAPI(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, Names)
}
