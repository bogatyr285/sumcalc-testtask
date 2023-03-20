package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"

	"github.com/bogatyr285/sumcalc-testtask/internal/mocks"
	jwtService "github.com/bogatyr285/sumcalc-testtask/internal/services/jwt"
)

func TestAuthOnly(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	tests := []struct {
		name             string
		mock             func(m *mocks.MockTokenVerifier)
		authorizationHdr string
		expectedStatus   int
	}{
		{
			name: "Valid token",
			mock: func(m *mocks.MockTokenVerifier) {
				m.EXPECT().VerifyToken([]byte("valid token")).Return(jwt.New(), nil)
			},
			authorizationHdr: "valid token",
			expectedStatus:   http.StatusOK,
		},
		{
			name: "Invalid token",
			mock: func(m *mocks.MockTokenVerifier) {
				m.EXPECT().VerifyToken([]byte("invalid token")).Return(nil, jwtService.ErrValidation)
			},
			authorizationHdr: "invalid token",
			expectedStatus:   http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tokenVerifierMock := mocks.NewMockTokenVerifier(ctrl)
			if tc.mock != nil {
				tc.mock(tokenVerifierMock)
			}

			router := gin.New()
			mw := AuthMiddleware{tokenVerifier: tokenVerifierMock}
			router.Use(mw.AuthOnly())
			router.GET("/", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Add("Authorization", tc.authorizationHdr)
			router.ServeHTTP(w, req)

			assert.Equal(t, w.Result().StatusCode, tc.expectedStatus)

		})
	}
}
