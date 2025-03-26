package oauth

import (
	"fmt"
	"go-image-processor/internal/db"
	"go-image-processor/internal/db/users"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/yandex"
)

var store *sessions.CookieStore
var AvailableProviders = []string{"yandex"}

const maxAge = 86400 * 30 // 30 days

func Init(r *gin.RouterGroup) {
	signingSecret := "fdsfsdfsdfdsfsdfsdfsd"
	store = sessions.NewCookieStore([]byte(signingSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   false,
		Domain:   "localhost",
	}

	gothic.Store = store

	yandexClientKey := os.Getenv("YANDEX_CLIENT_ID")
	yandexSecret := os.Getenv("YANDEX_SECRET")
	if yandexClientKey == "" || yandexSecret == "" {
		panic("oauth-keys are undefined")
	}

	goth.UseProviders(
		yandex.New(yandexClientKey, yandexSecret, "http://localhost:3000/oauth/yandex/callback"),
	)

	setupRoutes(r)
}

func setupRoutes(r *gin.RouterGroup) {
	oauth := r.Group("/oauth")
	oauth.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	oauth.GET("/:provider/callback", Callback)
	oauth.GET("/:provider/login", Login)
	oauth.GET("/logout", Logout)
	oauth.GET("/me", AuthMiddleware, GetCurrentUser)
}

// RetrieveUserBySession берёт session-куку, вытаскивает из неё данные по сессии
// и возвращает соответствующего пользователя из бд
func RetrieveUserBySession(c *gin.Context) (sessionData *db.User, err error) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		return nil, err
	}

	data, ok := session.Values["user"]
	if !ok {
		return nil, fmt.Errorf("invalid session")
	}

	sessionUser := data.(goth.User)

	user, exists := users.FindByEmail(&sessionUser.Email)
	if !exists {
		return nil, fmt.Errorf("user does not exist")
	}

	return user, nil
}

// AuthMiddleware can be used to protect auth-only routes
func AuthMiddleware(c *gin.Context) {
	account, err := RetrieveUserBySession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", account)
	c.Next()
}

func GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}
