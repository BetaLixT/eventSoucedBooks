package middlewares

import (
	"strings"

	"eventSourcedBooks/pkg/domain/auth"
	"eventSourcedBooks/pkg/infra/standard"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(svc *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := c.Request.Header["Authorization"]
		if !ok || len(token) < 1 {
			c.Error(standard.NewTokenMissingError())
			c.Abort()
			return
		}
		split := strings.Split(token[0], " ")
		if len(split) != 2 {
			c.Error(standard.NewTokenFormatInvalidError())
			c.Abort()
			return
		}

		claims, err := svc.ValidateToken(c, split[1])
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		c.Set(standard.APP_CLAIMS_KEY, claims)
		c.Next()
	}
}
