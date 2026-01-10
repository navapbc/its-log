package serve

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Read(c *gin.Context) {
	date := c.Param("date")
	operation := c.Param("operation")

	log.Println("READ " + date + " " + operation)

	// Everything worked.
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
