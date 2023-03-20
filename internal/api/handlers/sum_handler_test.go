package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/bogatyr285/sumcalc-testtask/internal/api/handlers"
	"github.com/bogatyr285/sumcalc-testtask/internal/mocks"
)

func TestSumHandler(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	type mocked struct {
		sumCalc *mocks.MockSumCalculator
		hasher  *mocks.MockHasher
	}
	tests := []struct {
		name         string
		request      string
		mocks        func(m *mocked)
		expectedCode int
		expectedBody string
	}{
		{
			name:    "valid request. zero response",
			request: `[]`,
			mocks: func(m *mocked) {
				m.hasher.EXPECT().Hash(int64(0)).Return("some_hash")
				m.sumCalc.EXPECT().SumNumbers(context.TODO(), []interface{}{}).Return(0)
			},
			expectedCode: http.StatusOK,
			expectedBody: marshal(&handlers.SumResponse{Sum: 0, Hash: "some_hash"}),
		},
		{
			name:    "invalid request",
			request: `[`,
			mocks: func(m *mocked) {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: marshal(handlers.ErrorResponse{Error: io.ErrUnexpectedEOF.Error()}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSumCalculator := mocks.NewMockSumCalculator(ctrl)
			mockHasher := mocks.NewMockHasher(ctrl)
			if tc.mocks != nil {
				tc.mocks(&mocked{
					sumCalc: mockSumCalculator,
					hasher:  mockHasher,
				})
			}

			handler := handlers.NewSumHandler(mockSumCalculator, mockHasher, zap.NewNop())

			w := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(w)
			ginCtx.Request, _ = http.NewRequest(
				"POST",
				"/sum",
				bytes.NewBuffer([]byte(tc.request)),
			)

			handler.SumHandler(ginCtx)

			b, _ := io.ReadAll(w.Body)
			require.Equal(t, w.Result().StatusCode, tc.expectedCode)
			require.Equal(t, string(b), tc.expectedBody)
		})
	}
}

func marshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
