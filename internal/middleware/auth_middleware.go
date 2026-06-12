package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wu-clan/lykn/internal/common"
	"github.com/wu-clan/lykn/pkg/response"
)

const userProfileKey = "user_profile"

var authSecret string

type UserProfile struct {
	ID       uint
	Username string
}

func SetAuthSecret(secret string) {
	authSecret = secret
}

func RequireLogin() gin.HandlerFunc {
	return authRequired(authSecret)
}

func AuthRequired(secret string) gin.HandlerFunc {
	return authRequired(secret)
}

func CurrentUserProfile(c *gin.Context) (*UserProfile, bool) {
	profile, ok := c.Get(userProfileKey)
	if !ok {
		return nil, false
	}
	typedProfile, ok := profile.(*UserProfile)
	return typedProfile, ok
}

func authRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret == "" {
			response.Error(c, http.StatusInternalServerError, "auth secret not configured")
			c.Abort()
			return
		}
		token, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			response.Error(c, http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}
		claims, err := common.ParseToken(token, secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		profile := &UserProfile{ID: claims.UserID, Username: claims.Username}
		ctx := common.WithUserContext(c.Request.Context(), profile.ID, profile.Username)
		c.Request = c.Request.WithContext(ctx)
		c.Set(userProfileKey, profile)
		c.Next()
	}
}

func bearerToken(header string) (string, bool) {
	if !strings.HasPrefix(header, "Bearer ") {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	return token, token != ""
}
