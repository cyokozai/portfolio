package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	NewRouter(testAssets(), newReadiness()).
		ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/profile", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var got Profile
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got.Name == "" {
		t.Error("name が空")
	}
	if got.Alias == "" {
		t.Error("alias が空")
	}
	if got.Label == "" {
		t.Error("label が空")
	}
	if len(got.Affiliations) == 0 {
		t.Error("affiliations が空")
	}
	for i, a := range got.Affiliations {
		if a == "" {
			t.Errorf("affiliations[%d] が空文字（プレースホルダを残さないこと）", i)
		}
	}
}
