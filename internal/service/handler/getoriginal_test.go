package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
	"github.com/go-chi/chi/v5"
)

func TestGetOriginal(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetOriginal)

	ts := httptest.NewServer(r)
	defer ts.Close()

	testShort := "testID"
	testURL := "https://example.com"
	cache.Cache[testShort] = testURL

	t.Run("Valid ID", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/" + testShort)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()
	})

	t.Run("Invalid ID", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/nonexistent")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/"+testShort, nil)
		responseRecorder := httptest.NewRecorder()

		GetOriginal(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Error code is not BadRequest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		GetOriginal(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
