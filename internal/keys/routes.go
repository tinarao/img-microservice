package keys

import (
	"go-image-processor/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetKeys(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := u.(*db.User)
	c.JSON(http.StatusOK, gin.H{
		"publicKey": user.PublicApiKey,
	})
}
