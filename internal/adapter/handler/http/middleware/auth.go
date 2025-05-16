package middleware

import (
	"strings"

	"clean-arch/internal/core/port"

	"github.com/gin-gonic/gin"
)

const (
	// authorizationHeaderKey is the key for authorization header in the request
	authorizationHeaderKey = "authorization"
	// authorizationType is the accepted authorization type
	authorizationType = "bearer"
	// authorizationPayloadKey is the key for authorization payload in the context
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware is a middleware to check if the user is authenticated
func AuthMiddleware(ts port.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.Request.Header.Get("Authorization")

		if len(authorizationHeader) == 0 {
			// err := domain.ErrEmptyAuthorizationHeader
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "Authorization header is missing",
			})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			// err := domain.ErrInvalidAuthorizationHeader
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid authorization header format",
			})
			return
		}

		currentAuthorizationType := strings.ToLower(fields[0])
		if currentAuthorizationType != authorizationType {
			// err := domain.ErrInvalidAuthorizationType
			// handleAbort(ctx, err)
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid authorization type",
			})
			return
		}

		accessToken := fields[1]
		payload, err := ts.VerifyToken(accessToken)
		if err != nil {
			// handleAbort(ctx, err)
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid token",
			})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
