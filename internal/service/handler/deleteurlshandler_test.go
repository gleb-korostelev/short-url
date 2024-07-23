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
	"github.com/gleb-korostelev/short-url/internal/service/handler"
	"github.com/gleb-korostelev/short-url/internal/worker"
	mock_db "github.com/gleb-korostelev/short-url/mocks"
)

func TestDeleteURLsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStore, workerPool)

	// Тестовые данные
	testURLs := []string{"http://example.com", "http://test.com"}
	jsonBody, _ := json.Marshal(testURLs)
	userID := "test-user-id"

	// Настройка контекста
	ctx := context.WithValue(context.Background(), config.UserContextKey, userID)

	tests := []struct {
		name           string
		body           *bytes.Buffer
		context        context.Context
		expectedStatus int
		setupMocks     func()
	}{
		{
			name:           "Successful Deletion",
			body:           bytes.NewBuffer(jsonBody),
			context:        ctx,
			expectedStatus: http.StatusAccepted,
			setupMocks: func() {
				mockStore.EXPECT().MarkURLsAsDeleted(gomock.Any(), userID, testURLs).Return(nil)
			},
		},
		{
			name:           "Unauthorized Access",
			body:           bytes.NewBuffer(jsonBody),
			context:        context.Background(), // Нет userID в контексте
			expectedStatus: http.StatusUnauthorized,
			setupMocks:     func() {},
		},
		{
			name:           "Invalid JSON Body",
			body:           bytes.NewBuffer([]byte(`invalid json`)),
			context:        ctx,
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/delete", tt.body)
			req = req.WithContext(tt.context)
			rr := httptest.NewRecorder()

			tt.setupMocks()

			svc.DeleteURLsHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
