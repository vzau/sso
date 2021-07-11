package v1

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
)

type certReturn struct {
	message string
	keys    jwk.Set
}

func GetCerts(c *gin.Context) {
	jkeyset := os.Getenv("SSO_JWKS")
	keyset, err := jwk.Parse([]byte(jkeyset))
	if err != nil {
		log4g.Category("controllers/certs").Error("Error parsing JWKs: " + err.Error())
		log4g.Category("controllers/certs").Debug(jkeyset)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse JWKs"})
		return
	}

	pub, err := jwk.PublicSetOf(keyset)
	if err != nil {
		log4g.Category("controllers/certs").Error("Error generating public keyset: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate public keyset"})
		return
	}
	jpub, _ := json.Marshal(pub)
	out := map[string]interface{}{}
	json.Unmarshal([]byte(jpub), &out)
	out["message"] = "OK"
	c.JSON(http.StatusOK, out)
}
