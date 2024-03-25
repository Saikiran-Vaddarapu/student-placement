package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func TestAuthorize(t *testing.T) {
	tests := []struct {
		apikey      string
		contentType string
		rescode     int
	}{
		{apikey: "a601e44e306e430f8dde987f65844f05", contentType: "application/json", rescode: 200},
		{apikey: "84dcb7c09b4a4af8a67f4577ffe9b255", contentType: "application/json", rescode: 200},
		{apikey: "84dcb7c09b4a4af8a67f4577ffe9b259", contentType: "application/json", rescode: 400},
		{apikey: "84dcb7c09b4a4af8a67f4577ffe9b255", contentType: "text", rescode: 415},
	}
	for _, val := range tests {
		req, _ := http.NewRequest("GET", "localhost:8080/", http.NoBody)
		req.Header.Set("X-API-KEY", val.apikey)
		req.Header.Set("Content-Type", val.contentType)

		res := httptest.NewRecorder()

		authorize(http.HandlerFunc(testHandler)).ServeHTTP(res, req)

		if res.Code != val.rescode {
			t.Errorf("Expected : %v Got : %v", val.rescode, res.Code)
		}
	}
}
