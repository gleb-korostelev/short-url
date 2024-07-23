package handler_test

import (
	"bytes"
	"context"
	"errors"
	"io"
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

func TestPostShorter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStore, workerPool)

	tests := []struct {
		name           string
		method         string
		userID         string
		body           io.Reader
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid Method",
			method:         "GET",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Only POST method is allowed\n",
		},
		{
			name:           "Unauthorized Access",
			method:         "POST",
			body:           bytes.NewReader([]byte("http://example.com")),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Error Reading Body",
			method:         "POST",
			body:           errReader{errors.New("mock read error")},
			userID:         "valid-user-id",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error reading request body\n",
		},
		{
			name:   "Successful Shorten",
			method: "POST",
			body:   bytes.NewReader([]byte("http://example.com")),
			userID: "valid-user-id",
			setupMocks: func() {
				mockStore.EXPECT().SaveUniqueURL(gomock.Any(), "http://example.com", "valid-user-id").Return("http://short.url", http.StatusCreated, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://short.url",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, "/", tc.body)
			if tc.userID != "" {
				ctx := context.WithValue(req.Context(), config.UserContextKey, tc.userID)
				req = req.WithContext(ctx)
			}
			rr := httptest.NewRecorder()

			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			svc.PostShorter(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedBody != "" {
				responseBody := rr.Body.String()
				assert.Equal(t, tc.expectedBody, responseBody)
			}
		})
	}
}

// errReader helps simulate an error while reading the request body.
type errReader struct {
	err error
}

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}
