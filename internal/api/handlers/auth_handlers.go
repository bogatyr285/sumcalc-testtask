//go:generate mockgen -source auth_handlers.go  -destination=../../mocks/mock_auth_handlers.go -package=mocks
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TokenIssuer interface {
	IssueToken(userID string) (string, error)
}

type AuthHandlers struct {
	tokenIssuer TokenIssuer
	logger      *zap.Logger
}

func NewAuthHandlers(tokenIssuer TokenIssuer, logger *zap.Logger) *AuthHandlers {
	return &AuthHandlers{
		tokenIssuer: tokenIssuer,
		logger:      logger.Named("auth-handler"),
	}
}

func (h *AuthHandlers) AuthHandler(c *gin.Context) {
	authReq := &AuthRequest{}
	if err := c.ShouldBindJSON(authReq); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Error: err.Error()})
		return
	}
	if err := authReq.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, &ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.tokenIssuer.IssueToken(authReq.Username)
	if err != nil {
		// dont show to user actual error, just log it
		h.logger.Error("issuing token err", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorResponse{Error: "internal error"})
		return
	}
	resp := &AuthResponse{Payload: token}
	c.JSON(http.StatusOK, resp)
}
