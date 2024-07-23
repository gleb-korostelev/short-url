package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/handler"
	"github.com/gleb-korostelev/short-url.git/internal/service/utils"
	"github.com/gleb-korostelev/short-url.git/internal/storage/inmemory"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
	mock_db "github.com/gleb-korostelev/short-url.git/mocks"
	"github.com/golang/mock/gomock"
)

func ExampleAPIService_GetOriginal() {
	store := inmemory.NewMemoryStorage(cache.Cache)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	apiService := handler.NewAPIService(store, workerPool)
	server := httptest.NewServer(http.HandlerFunc(apiService.GetOriginal))
	defer server.Close()

	resp, err := http.Get(server.URL + "/get_original?id=123")
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
}

func ExampleAPIService_PostShorter() {
	store := inmemory.NewMemoryStorage(cache.Cache)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	apiService := handler.NewAPIService(store, workerPool)
	server := httptest.NewServer(http.HandlerFunc(apiService.PostShorter))
	defer server.Close()

	resp, err := http.Post(server.URL+"/post_shorter", "application/text", strings.NewReader("http://example.com"))
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
}
func ExampleAPIService_DeleteURLsHandler() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	mockWorkerPool := worker.NewDBWorkerPool(1)
	defer mockWorkerPool.Shutdown()

	mockStore.EXPECT().MarkURLsAsDeleted(gomock.Any(), "user-123", []string{"short1", "short2"}).Return(nil)

	svc := handler.NewAPIService(mockStore, mockWorkerPool)

	urlsToDelete := []string{"short1", "short2"}
	body, _ := json.Marshal(urlsToDelete)
	req := httptest.NewRequest("DELETE", "/delete-urls", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), config.UserContextKey, "user-123"))

	rr := httptest.NewRecorder()

	svc.DeleteURLsHandler(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		panic("expected HTTP 202 Accepted")
	}

	print("HTTP Status:", rr.Code)
}

func ExampleAPIService_GetUserURLs() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)

	expectedURLs := []models.UserURLs{
		{ShortURL: "http://short.url", OriginalURL: "http://original.url"},
	}
	mockStore.EXPECT().GetAllURLS(context.Background(), "user-123", config.BaseURL).Return(expectedURLs, nil)

	svc := handler.NewAPIService(mockStore, workerPool)

	req := httptest.NewRequest("GET", "/user/urls", nil)

	rr := httptest.NewRecorder()
	utils.SetJWTInCookie(rr, "user-123")

	svc.GetUserURLs(rr, req)

	response := rr.Result()
	defer response.Body.Close()

	var urls []models.UserURLs
	if err := json.NewDecoder(response.Body).Decode(&urls); err != nil {
		panic(err)
	}

	fmt.Println("HTTP Status:", response.StatusCode)
	for _, url := range urls {
		fmt.Printf("Short URL: %s, Original URL: %s\n", url.ShortURL, url.OriginalURL)
	}
}

func ExampleAPIService_Ping() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStore, workerPool)

	mockStore.EXPECT().Ping(context.Background()).Return(http.StatusOK, nil)

	req := httptest.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()

	svc.Ping(rr, req)

	response := rr.Result()
	defer response.Body.Close()

	fmt.Println("HTTP Status:", response.StatusCode)
}
