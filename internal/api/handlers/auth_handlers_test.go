package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/bogatyr285/sumcalc-testtask/internal/api"
	"github.com/bogatyr285/sumcalc-testtask/internal/api/handlers"
	"github.com/bogatyr285/sumcalc-testtask/internal/mocks"
)

func TestAuthHandler(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	tests := []struct {
		name           string
		authRequest    handlers.AuthRequest
		mock           func(m *mocks.MockTokenIssuer)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "valid request",
			authRequest: handlers.AuthRequest{Username: "Bob", Password: "1337"},
			mock: func(m *mocks.MockTokenIssuer) {
				m.EXPECT().IssueToken("Bob").Return("token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"payload":"token"}`,
		},
		{
			name:        "invalid request. empty username",
			authRequest: handlers.AuthRequest{Username: "", Password: "1337"},
			mock: func(m *mocks.MockTokenIssuer) {
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":"Key: 'AuthRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`,
		},
		{
			name:        "invalid request. empty password",
			authRequest: handlers.AuthRequest{Username: "Trudy", Password: ""},
			mock: func(m *mocks.MockTokenIssuer) {
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":"Key: 'AuthRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:        "token issuing error",
			authRequest: handlers.AuthRequest{Username: "Alice", Password: "1337"},
			mock: func(m *mocks.MockTokenIssuer) {
				m.EXPECT().IssueToken("Alice").Return("", errors.New("some rare internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   fmt.Sprintf(`{"error":"%v"}`, api.ErrInternal.Error()),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockTokenIssuer := mocks.NewMockTokenIssuer(ctrl)
			if tc.mock != nil {
				tc.mock(mockTokenIssuer)
			}

			handler := handlers.NewAuthHandlers(mockTokenIssuer, zap.NewNop())

			w := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(w)
			jsonValue, _ := json.Marshal(tc.authRequest)
			ginCtx.Request, _ = http.NewRequest(
				"POST",
				"/auth",
				bytes.NewBuffer(jsonValue),
			)

			handler.AuthHandler(ginCtx)

			b, _ := ioutil.ReadAll(w.Body)
			require.Equal(t, w.Result().StatusCode, tc.expectedStatus)
			require.Equal(t, string(b), tc.expectedBody)
		})
	}
}
