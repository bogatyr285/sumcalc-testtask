//go:generate mockgen -source sum_handler.go  -destination=../../mocks/mock_sum_handler.go -package=mocks
package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SumCalculator - find all digits in json query and sum it
type SumCalculator interface {
	SumNumbers(ctx context.Context, v interface{}) int
}

// Hasher - hashes number with some hash alg
type Hasher interface {
	Hash(n int64) string
}

type SumHandler struct {
	sumCalculator SumCalculator
	hasher        Hasher
	logger        *zap.Logger
}

func NewSumHandler(sumCalculator SumCalculator, hasher Hasher, logger *zap.Logger) *SumHandler {
	return &SumHandler{
		sumCalculator: sumCalculator,
		hasher:        hasher,
		logger:        logger.Named("sum-handler"),
	}
}

func (h *SumHandler) SumHandler(c *gin.Context) {
	var v interface{}

	if err := json.NewDecoder(c.Request.Body).Decode(&v); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorResponse{Error: err.Error()})
		return
	}
	// todo we must have tracing middleware and pass uniq traceid into service func inside ctx
	sum := h.sumCalculator.SumNumbers(c.Request.Context(), v)
	// todo we can make separate module which will do "sum & hashing" and here make just one func call
	// but it's another layer of abstraction which's not necessary right now, but can be considered on app grow
	hashedSum := h.hasher.Hash(int64(sum))
	c.JSON(http.StatusOK, &SumResponse{
		Sum:  sum,
		Hash: hashedSum,
	})
}
