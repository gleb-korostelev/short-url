package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/service/handler"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
	mock_db "github.com/gleb-korostelev/short-url.git/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPostShorterJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStore, workerPool)

	tests := []struct {
		name           string
		method         string
		userID         string
		requestBody    string
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid Method",
			method:         "GET",
			requestBody:    "{}",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Only POST method is allowed\n",
		},
		{
			name:           "Unauthorized Access",
			method:         "POST",
			requestBody:    "{}",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid JSON",
			method:         "POST",
			userID:         "valid-user-id",
			requestBody:    "{invalid_json}",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Database Error",
			method:      "POST",
			userID:      "valid-user-id",
			requestBody: `{"url":"http://example.com"}`,
			setupMocks: func() {
				mockStore.EXPECT().
					SaveUniqueURL(gomock.Any(), "http://example.com", "valid-user-id").
					Return("", http.StatusInternalServerError, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Error with saving\n",
		},
		{
			name:        "Successful Shorten",
			method:      "POST",
			userID:      "valid-user-id",
			requestBody: `{"url":"http://example.com"}`,
			setupMocks: func() {
				mockStore.EXPECT().
					SaveUniqueURL(gomock.Any(), "http://example.com", "valid-user-id").
					Return("http://short.url", http.StatusCreated, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"result":"http://short.url"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, "/shorten-json", bytes.NewBufferString(tc.requestBody))
			if tc.userID != "" {
				ctx := context.WithValue(req.Context(), config.UserContextKey, tc.userID)
				req = req.WithContext(ctx)
			}
			rr := httptest.NewRecorder()

			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			svc.PostShorterJSON(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedBody != "" {
				responseBody := rr.Body.String()
				if tc.expectedStatus == http.StatusBadRequest || tc.expectedStatus == http.StatusInternalServerError {
					assert.Equal(t, tc.expectedBody, responseBody)
				} else {
					assert.JSONEq(t, tc.expectedBody, responseBody)
				}
			}
		})
	}
}
