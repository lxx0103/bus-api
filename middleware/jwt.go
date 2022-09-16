package middleware

import (
	"bus-api/core/response"
	"bus-api/service"
	"errors"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) <= len(BEARER_SCHEMA) {
			response.ResponseUnauthorized(c, "AuthError", errors.New("NO AUTH HEADER"))
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		if tokenString == "" {
			response.ResponseUnauthorized(c, "AuthError", errors.New("JWT AUTH ERROR"))
			return
		}
		claims, err := service.JWTAuthService().ParseToken(tokenString)
		if err != nil {
			response.ResponseUnauthorized(c, "AuthError", errors.New("JWT AUTH ERROR"))
			return
		}
		// var claims service.CustomClaims
		// claims.UserID = 1
		// claims.Username = "lewis"
		c.Set("claims", claims)
		c.Next()
	}
}
