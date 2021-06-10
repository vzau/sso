package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, message string) {
	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": message})
	c.Abort()
}
