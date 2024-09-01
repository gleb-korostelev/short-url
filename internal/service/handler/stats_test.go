package handler_test

import (
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

func TestStatsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStore, workerPool)

	tests := []struct {
		name           string
		trustedSubnet  string
		clientIP       string
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "Allowed IP within subnet",
			trustedSubnet: "192.168.1.0/24",
			clientIP:      "192.168.1.5",
			setupMocks: func() {
				mockStore.EXPECT().GetStats(gomock.Any()).Return(150, 25, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"urls":150,"users":25}`,
		},
		{
			name:           "Disallowed IP outside subnet",
			trustedSubnet:  "192.168.1.0/24",
			clientIP:       "192.168.2.5",
			setupMocks:     func() {},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Forbidden\n",
		},
		{
			name:           "Empty subnet configuration",
			trustedSubnet:  "",
			clientIP:       "192.168.1.5",
			setupMocks:     func() {},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Forbidden\n",
		},
		{
			name:          "Storage error",
			trustedSubnet: "192.168.1.0/24",
			clientIP:      "192.168.1.5",
			setupMocks: func() {
				mockStore.EXPECT().GetStats(gomock.Any()).Return(0, 0, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal server error\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the request and response recorder
			req, err := http.NewRequest("GET", "/api/internal/stats", nil)
			assert.NoError(t, err)
			req.Header.Set("X-Real-IP", tc.clientIP)
			rr := httptest.NewRecorder()

			tc.setupMocks()
			// Setup the handler
			config.TrustedSubnet = tc.trustedSubnet
			svc.StatsHandler(rr, req)

			// Check the status code and response body
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if rr.Code == http.StatusOK {
				assert.JSONEq(t, tc.expectedBody, rr.Body.String())
			} else {
				assert.Equal(t, tc.expectedBody, rr.Body.String())
			}
		})
	}
}
