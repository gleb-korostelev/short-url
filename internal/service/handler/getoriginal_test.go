package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gleb-korostelev/short-url/internal/config"
	"github.com/gleb-korostelev/short-url/internal/service/handler"
	"github.com/gleb-korostelev/short-url/internal/storage/repository"
	"github.com/gleb-korostelev/short-url/internal/worker"
	mock_db "github.com/gleb-korostelev/short-url/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func BenchmarkProcessURLs(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	mockdb := mock_db.NewMockDatabaseI(ctrl)
	store := repository.NewDBStorage(mockdb)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)

	svc := handler.NewAPIService(store, workerPool)

	testShort := "testID"

	request, _ := http.NewRequest(http.MethodPost, "/"+testShort, nil)
	responseRecorder := httptest.NewRecorder()

	svc.GetOriginal(responseRecorder, request)
}

// MockStore - mock for Storage interface.
type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetOriginalLink(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func TestGetOriginal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStorage, workerPool)

	r := chi.NewRouter()
	r.Get("/{id}", svc.GetOriginal)

	tests := []struct {
		name         string
		id           string
		mockResponse string
		mockError    error
		expectedCode int
		expectedLoc  string
	}{
		{
			name:         "Valid ID",
			id:           "123",
			mockResponse: "http://original.url/example",
			mockError:    nil,
			expectedCode: http.StatusTemporaryRedirect,
			expectedLoc:  "http://original.url/example",
		},
		{
			name:         "URL Gone",
			id:           "gone",
			mockResponse: "",
			mockError:    config.ErrGone,
			expectedCode: http.StatusGone,
			expectedLoc:  "",
		},
		{
			name:         "Invalid ID",
			id:           "invalid",
			mockResponse: "",
			mockError:    errors.New("not found"),
			expectedCode: http.StatusBadRequest,
			expectedLoc:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/"+tt.id, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			mockStorage.EXPECT().
				GetOriginalLink(gomock.Any(), tt.id).
				Return(tt.mockResponse, tt.mockError).
				Times(1)

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedLoc != "" {
				assert.Equal(t, tt.expectedLoc, rr.Header().Get("Location"))
			}
		})
	}
}
