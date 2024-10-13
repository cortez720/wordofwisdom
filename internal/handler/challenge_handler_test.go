package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mocks "github.com/cortez720/wordofwisdom/internal/handler/mocks"
)

func TestPowHandler_Challenge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Pow interface
	mockPow := mocks.NewMockpow(ctrl)

	// Set up expected behavior
	challengeBytes := []byte("mockChallenge")
	mockPow.EXPECT().Challenge().Return(challengeBytes)

	// Create an instance of PowHandler with the mock
	handler := NewPow(mockPow, nil)

	// Create a request and a response recorder
	req := httptest.NewRequest(http.MethodGet, "/challenge", nil)
	rr := httptest.NewRecorder()

	// Call the Challenge method
	handler.Challenge(rr, req)

	// Assertions
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, challengeBytes, rr.Body.Bytes())
}

func TestPowHandler_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPow := mocks.NewMockpow(ctrl)
	mockSvc := mocks.NewMockquoteService(ctrl)

	handler := NewPow(mockPow, mockSvc)

	t.Run("valid challenge and solution", func(t *testing.T) {
		challenge := []byte("validChallenge")
		solution := []byte("validSolution")

		mockPow.EXPECT().Verify(challenge, solution).Return(nil)
		mockSvc.EXPECT().GetWordOfWisdom(gomock.Any()).Return("Keep pushing!", nil)

		req := httptest.NewRequest(http.MethodPost, "/validate", nil)
		req.Form = url.Values{
			challengeArg: {string(challenge)},
			solutionArg:  {string(solution)},
		}
		rr := httptest.NewRecorder()

		handler.Validate(rr, req)

		// Assertions
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "Quote of the day: Keep pushing!", rr.Body.String())
	})

	t.Run("invalid PoW solution", func(t *testing.T) {
		challenge := []byte("invalidChallenge")
		solution := []byte("invalidSolution")

		testError := errors.New("invalid solution error")

		mockPow.EXPECT().Verify(challenge, solution).Return(testError)

		req := httptest.NewRequest(http.MethodPost, "/validate", nil)
		req.Form = url.Values{
			challengeArg: {string(challenge)},
			solutionArg:  {string(solution)},
		}
		rr := httptest.NewRecorder()

		handler.Validate(rr, req)

		// Assertions
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		require.Equal(t, "Invalid PoW solution\n", rr.Body.String())
	})

	t.Run("quote service error", func(t *testing.T) {
		challenge := []byte("validChallenge")
		solution := []byte("validSolution")

		mockPow.EXPECT().Verify(challenge, solution).Return(nil)
		mockSvc.EXPECT().GetWordOfWisdom(gomock.Any()).Return("", fmt.Errorf("service error"))

		req := httptest.NewRequest(http.MethodPost, "/validate", nil)
		req.Form = url.Values{
			challengeArg: {string(challenge)},
			solutionArg:  {string(solution)},
		}
		rr := httptest.NewRecorder()

		handler.Validate(rr, req)

		// Assertions
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		require.Equal(t, "Internal error.\n", rr.Body.String())
	})
}
