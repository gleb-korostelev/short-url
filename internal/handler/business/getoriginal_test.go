package business

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gleb-korostelev/short-url.git/internal/config"
)

func TestGetOriginal(t *testing.T) {
	testShort := "testID"
	testURL := "https://example.com"
	config.Cache[testShort] = testURL

	t.Run("Valid ID", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/"+testShort, nil)
		responseRecorder := httptest.NewRecorder()

		GetOriginal(responseRecorder, request)

		result := responseRecorder.Body.String()
		if result != testURL {
			t.Errorf("Expected %s, got %s", testURL, result)
		}
		if status := responseRecorder.Code; status != http.StatusTemporaryRedirect {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTemporaryRedirect)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/nonexistent", nil)
		responseRecorder := httptest.NewRecorder()

		GetOriginal(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusBadRequest {
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
}
