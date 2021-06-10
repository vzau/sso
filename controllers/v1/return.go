package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReturn(c *gin.Context) {
	token := c.Param("token")
	if len(token) < 1 {
		c.HTML(http.StatusInternalServerError, "error.tmpl", "Invalid response received from Authenticator or Authentication cancelled.")
		return
	}

	//	key := []byte(os.Getenv("VATUSA_ULS_KEY"))

}
