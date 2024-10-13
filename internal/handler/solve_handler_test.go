package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mocks "github.com/cortez720/wordofwisdom/internal/handler/mocks"
)

func TestSolverHandler_Solve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSolverSvc := mocks.NewMockSolverService(ctrl)
	handler := NewSolver(mockSolverSvc)

	t.Run("successful solve", func(t *testing.T) {
		expectedQuote := []byte("Keep going!")
		mockSolverSvc.EXPECT().Solve().Return(expectedQuote, nil)

		req := httptest.NewRequest(http.MethodGet, "/solve", nil)
		rr := httptest.NewRecorder()

		handler.Solve(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, string(expectedQuote), rr.Body.String())
	})

	t.Run("solve service error", func(t *testing.T) {
		testError := errors.New("invalid solution error")

		mockSolverSvc.EXPECT().Solve().Return(nil, testError)

		req := httptest.NewRequest(http.MethodGet, "/solve", nil)
		rr := httptest.NewRecorder()

		handler.Solve(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		require.Equal(t, "Internal error.\n", rr.Body.String())
	})
}
