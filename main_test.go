package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCatalogHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	w := httptest.NewRecorder()

	catalogHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}

	var videos []Video
	err := json.Unmarshal(w.Body.Bytes(), &videos)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(videos) == 0 {
		t.Error("Expected at least one video in the catalog, but got none")
	}
}
