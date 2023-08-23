package http

import (
	"net/http/httptest"
	"testing"
)

func Test_handleIndex_ReturnsCorrectStatus(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	NewServer().handleIndex()(w, r)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}
