//go:generate mockgen -source auth.go  -destination=../../mocks/mock_auth_middleware.go -package=mocks
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const AuthorizationHeader = "Authorization"

type TokenVerifier interface {
	VerifyToken(payload []byte) (jwt.Token, error)
}

type AuthMiddleware struct {
	tokenVerifier TokenVerifier
}

func NewAuthMiddleware(tokenVerifier TokenVerifier) *AuthMiddleware {
	return &AuthMiddleware{tokenVerifier}
}

// todo here we can add roles and also check it
func (a *AuthMiddleware) AuthOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader(AuthorizationHeader)
		_, err := a.tokenVerifier.VerifyToken([]byte(bearerToken))
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
