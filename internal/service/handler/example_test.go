package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/handler"
	"github.com/gleb-korostelev/short-url.git/internal/service/utils"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
	mock_db "github.com/gleb-korostelev/short-url.git/mocks"
	"github.com/golang/mock/gomock"
)

func ExampleAPIService_GetOriginal() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	apiService := handler.NewAPIService(mockStore, workerPool)
	server := httptest.NewServer(http.HandlerFunc(apiService.GetOriginal))
	defer server.Close()

	resp, err := http.Get(server.URL + "/{123}")
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
}

func ExampleAPIService_PostShorter() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	apiService := handler.NewAPIService(mockStore, workerPool)
	server := httptest.NewServer(http.HandlerFunc(apiService.PostShorter))
	defer server.Close()

	resp, err := http.Post(server.URL+"/", "application/text", strings.NewReader("http://example.com"))
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
	workerPool := worker.NewDBWorkerPool(1)
	defer workerPool.Shutdown()

	mockStore.EXPECT().MarkURLsAsDeleted(gomock.Any(), "user-123", []string{"short1", "short2"}).Return(nil)

	svc := handler.NewAPIService(mockStore, workerPool)

	urlsToDelete := []string{"short1", "short2"}
	body, _ := json.Marshal(urlsToDelete)
	req := httptest.NewRequest("DELETE", "/api/user/urls", bytes.NewBuffer(body))
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

	req := httptest.NewRequest("GET", "/api/user/urls", nil)

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

func ExampleAPIService_PostShorterJSON() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStore, workerPool)

	payload := models.URLPayload{
		URL: "https://example.com",
	}
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(context.Background(), config.UserContextKey, "test-user-id"))

	rr := httptest.NewRecorder()

	svc.PostShorterJSON(rr, req)

	fmt.Printf("Status Code: %d\n", rr.Code)
	fmt.Printf("Body: %s\n", rr.Body.String())
}

func ExampleAPIService_ShortenBatchHandler() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStore := mock_db.NewMockStorage(ctrl)
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	svc := handler.NewAPIService(mockStore, workerPool)

	batchRequest := []models.ShortenBatchRequestItem{
		{CorrelationID: "1", OriginalURL: "https://example.com"},
		{CorrelationID: "2", OriginalURL: "https://example.org"},
	}
	jsonBody, _ := json.Marshal(batchRequest)

	req, _ := http.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(context.Background(), config.UserContextKey, "test-user-id"))

	rr := httptest.NewRecorder()

	svc.ShortenBatchHandler(rr, req)

	fmt.Printf("Status Code: %d\n", rr.Code)
	fmt.Printf("Body: %s\n", rr.Body.String())
}
