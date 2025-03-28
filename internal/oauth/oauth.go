package oauth

import (
	"fmt"
	"go-image-processor/internal/db"
	"go-image-processor/internal/db/users"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/yandex"
	"github.com/michaeljs1990/sqlitestore"
)

var store *sqlitestore.SqliteStore
var AvailableProviders = []string{"yandex"}

const maxAge = 86400 * 30 // 30 days

func Init(r *gin.RouterGroup) {
	signingSecret := "fdsfsdfsdfdsfsdfsdfsd"

	var err error
	store, err = sqlitestore.NewSqliteStore("test.db", "sessions", "/", maxAge, []byte(signingSecret))
	if err != nil {
		panic(err)
	}

	gothic.Store = store

	yandexClientKey := os.Getenv("YANDEX_CLIENT_ID")
	yandexSecret := os.Getenv("YANDEX_SECRET")
	if yandexClientKey == "" || yandexSecret == "" {
		panic("oauth-keys are undefined")
	}

	goth.UseProviders(
		yandex.New(yandexClientKey, yandexSecret, "http://localhost:3000/api/oauth/yandex/callback"),
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

func AuthMiddleware(c *gin.Context) {
	account, err := RetrieveUserBySession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", account)
	c.Next()
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

func GetCurrentUser(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	user := u.(*db.User)
	c.JSON(http.StatusOK, user)
}
