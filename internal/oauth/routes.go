package oauth

import (
	"go-image-processor/internal/db/users"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func Login(c *gin.Context) {
	provider := c.Param("provider")
	if !slices.Contains(AvailableProviders, provider) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)

	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func Logout(c *gin.Context) {
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		slog.Error("failed to logout using gothic", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Callback(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		slog.Error("failed authorization attempt", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to authorize"})
		return
	}

	_, err = users.CompleteAuthorization(&user)
	if err != nil {
		slog.Error("failed to complete authorization", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to authorize"})
		return
	}

	session, _ := store.Get(c.Request, "session")
	session.Values["user"] = user
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		slog.Error("failed to save the session", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/api/accounts/p/me")
}
