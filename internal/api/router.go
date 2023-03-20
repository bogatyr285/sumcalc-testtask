package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/bogatyr285/sumcalc-testtask/internal/api/handlers"
	"github.com/bogatyr285/sumcalc-testtask/internal/api/middlewares"
	"github.com/bogatyr285/sumcalc-testtask/internal/buildinfo"
)

func NewRouter(
	tokenIssuer handlers.TokenIssuer,
	tokenVerifier middlewares.TokenVerifier,
	sumCalculator handlers.SumCalculator,
	hasher handlers.Hasher,
	logger *zap.Logger,
) *gin.Engine {
	r := newRouter()

	authMiddleware := middlewares.NewAuthMiddleware(tokenVerifier)
	authHandlers := handlers.NewAuthHandlers(tokenIssuer, logger)
	sumHandler := handlers.NewSumHandler(sumCalculator, hasher, logger)

	r.POST("/auth", authHandlers.AuthHandler)
	r.POST("/sum", authMiddleware.AuthOnly(), sumHandler.SumHandler)

	r.GET("/build", buildinfo.Handler(buildinfo.New()))

	return r
}

func newRouter() *gin.Engine {
	r := gin.New()
	// todo
	//r.Use(jaeger)
	//r.Use(advancedPanicStackTracing)
	//r.Use(NewLoggerMiddleware)

	return r
}
