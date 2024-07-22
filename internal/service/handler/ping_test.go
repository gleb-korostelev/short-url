package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_db "github.com/gleb-korostelev/short-url.git/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	svc := &APIService{store: mockStore}

	tests := []struct {
		name         string
		mockStatus   int
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful Ping",
			mockStatus:   http.StatusOK,
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Failed Ping with error",
			mockStatus:   http.StatusInternalServerError,
			mockError:    fmt.Errorf("database connection error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "database connection error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/ping", nil)
			rr := httptest.NewRecorder()

			mockStore.EXPECT().Ping(context.Background()).Return(tt.mockStatus, tt.mockError)

			svc.Ping(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				responseBody := rr.Body.String()
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}
