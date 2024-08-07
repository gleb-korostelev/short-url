package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/handler"
	"github.com/gleb-korostelev/short-url.git/internal/service/utils"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
	mock_db "github.com/gleb-korostelev/short-url.git/mocks"
)

func TestGetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStore, workerPool)

	tests := []struct {
		name            string
		userID          string
		mockSetup       func()
		expectedStatus  int
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{
			name:   "Successful retrieval of URLs",
			userID: "valid-user-id",
			mockSetup: func() {
				mockStore.EXPECT().
					GetAllURLS(gomock.Any(), "valid-user-id", config.BaseURL).
					Return([]models.UserURLs{{ShortURL: "http://short.url", OriginalURL: "http://original.url"}}, nil)
			},
			expectedStatus:  http.StatusOK,
			expectedBody:    `[{"short_url":"http:\/\/short.url","original_url":"http:\/\/original.url"}]`,
			expectedHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:   "Internal server error on store failure",
			userID: "valid-user-id",
			mockSetup: func() {
				mockStore.EXPECT().
					GetAllURLS(gomock.Any(), "valid-user-id", config.BaseURL).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "No content when no URLs are found",
			userID: "valid-user-id",
			mockSetup: func() {
				mockStore.EXPECT().
					GetAllURLS(gomock.Any(), "valid-user-id", config.BaseURL).
					Return([]models.UserURLs{}, nil)
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/user/urls", nil)
			rr := httptest.NewRecorder()

			utils.SetJWTInCookie(rr, tc.userID)
			tc.mockSetup()

			response := rr.Result()
			defer response.Body.Close()
			cookies := response.Cookies()

			if len(cookies) > 0 {
				req.AddCookie(cookies[0])
			}
			svc.GetUserURLs(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedBody != "" {
				body := rr.Body.String()
				assert.JSONEq(t, tc.expectedBody, body)
			}
			for key, value := range tc.expectedHeaders {
				assert.Equal(t, value, rr.Header().Get(key))
			}
		})
	}
}
