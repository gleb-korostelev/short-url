package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/storage/repository"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
	mock_db "github.com/gleb-korostelev/short-url.git/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func BenchmarkProcessURLs(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	mockdb := mock_db.NewMockDatabaseI(ctrl)
	store := repository.NewDBStorage(mockdb)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)

	svc := NewAPIService(store, workerPool)

	testShort := "testID"

	request, _ := http.NewRequest(http.MethodPost, "/"+testShort, nil)
	responseRecorder := httptest.NewRecorder()

	svc.GetOriginal(responseRecorder, request)
}

func TestGetOriginal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_db.NewMockDatabaseI(ctrl)
	r := chi.NewRouter()
	store := repository.NewDBStorage(mockdb)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)

	svc := NewAPIService(store, workerPool)
	r.Get("/{id}", svc.GetOriginal)

	ts := httptest.NewServer(r)
	defer ts.Close()

	testShort := "testID"
	testURL := "https://example.com"
	var testdata models.URLData
	testdata.OriginalURL = testURL
	testdata.ShortURL = testShort
	testdata.UUID = uuid.New()
	cache.Cache[testURL] = testdata
	// cache.Cache = append(cache.Cache, testdata)

	t.Run("Unsupported Method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/"+testShort, nil)
		responseRecorder := httptest.NewRecorder()

		svc.GetOriginal(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Error code is not BadRequest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		svc.GetOriginal(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
