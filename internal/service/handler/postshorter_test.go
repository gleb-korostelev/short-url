package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gleb-korostelev/short-url.git/internal/service/business"
)

func TestPostShorter(t *testing.T) {
	business.MockCacheURL = func(originalURL string) string {
		return "http://short.url/test"
	}
	defer func() { business.MockCacheURL = nil }()

	t.Run("Create Short URL", func(t *testing.T) {
		originalURL := "https://example.com"
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(originalURL))
		responseRecorder := httptest.NewRecorder()

		PostShorter(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		expected := "http://short.url/test"
		if response := responseRecorder.Body.String(); response != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", response, expected)
		}
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		PostShorter(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Error code is not BadRequest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/test", nil)
		responseRecorder := httptest.NewRecorder()

		GetOriginal(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
