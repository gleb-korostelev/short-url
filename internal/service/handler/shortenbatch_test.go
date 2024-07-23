package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/gleb-korostelev/short-url/internal/config"
	"github.com/gleb-korostelev/short-url/internal/models"
	"github.com/gleb-korostelev/short-url/internal/service/handler"
	"github.com/gleb-korostelev/short-url/internal/worker"
	mock_db "github.com/gleb-korostelev/short-url/mocks"
)

func TestShortenBatchHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStore, workerPool)

	tests := []struct {
		name            string
		userID          string
		requestBody     []models.ShortenBatchRequestItem
		setupMocks      func()
		expectedStatus  int
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{
			name:   "Successful batch shorten",
			userID: "valid-user-id",
			requestBody: []models.ShortenBatchRequestItem{
				{CorrelationID: "1", OriginalURL: "http://example.com"},
			},
			setupMocks: func() {
				mockStore.EXPECT().SaveURL(gomock.Any(), "http://example.com", "valid-user-id").Return("http://short.url", nil)
			},
			expectedStatus:  http.StatusCreated,
			expectedBody:    `[{"correlation_id":"1","short_url":"http://short.url"}]`,
			expectedHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:           "Unauthorized without user context",
			userID:         "",
			requestBody:    []models.ShortenBatchRequestItem{{CorrelationID: "1", OriginalURL: "http://example.com"}},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid JSON body",
			userID:         "valid-user-id",
			requestBody:    nil,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Empty batch is not allowed\n",
		},
		{
			name:           "Empty batch request",
			userID:         "valid-user-id",
			requestBody:    []models.ShortenBatchRequestItem{},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Empty batch is not allowed\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/batch", bytes.NewBuffer(body))
			if tc.userID != "" {
				ctx := context.WithValue(req.Context(), config.UserContextKey, tc.userID)
				req = req.WithContext(ctx)
			}
			rr := httptest.NewRecorder()

			tc.setupMocks()

			svc.ShortenBatchHandler(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedBody != "" {
				responseBody := rr.Body.String()
				if tc.expectedStatus == http.StatusBadRequest {
					assert.Equal(t, tc.expectedBody, responseBody)
				} else {
					assert.JSONEq(t, tc.expectedBody, responseBody)
				}
			}
			for key, value := range tc.expectedHeaders {
				assert.Equal(t, value, rr.Header().Get(key))
			}
		})
	}
}
