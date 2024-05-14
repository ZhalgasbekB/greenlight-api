package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateMovie(t *testing.T) {
	app := newTestApp()
	router := app.routes()

	movieJSON := `{"title":"Inception","year":2000,"runtime":"48 mins","genres":["action","sci-fi"]}`
	req, err := http.NewRequest(http.MethodPost, "/v1/movies", strings.NewReader(movieJSON))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected HTTP status code 201 (Created)")
	assert.Contains(t, rr.Body.String(), "Inception", "Response body should contain 'Inception'")
}
func TestUpdateMovie(t *testing.T) {
	app := newTestApp()
	router := app.routes()

	updateJSON := `{"title":"Inception Updated","year":2010,"runtime":"149 mins","genres":["action","sci-fi"]}`
	req, err := http.NewRequest(http.MethodPatch, "/v1/movies/1", strings.NewReader(updateJSON))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP status code 200 OK")
	assert.Contains(t, rr.Body.String(), "Inception Updated", "Response body should contain updated title 'Inception Updated'")
}

func TestDeleteMovie(t *testing.T) {
	app := newTestApp()
	router := app.routes()

	req, err := http.NewRequest(http.MethodDelete, "/v1/movies/1", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP status code 200 OK")
	assert.Contains(t, rr.Body.String(), "", "Response body should be empty after deletion")
}

func TestListMovies(t *testing.T) {
	app := newTestApp()
	router := app.routes()

	req, err := http.NewRequest(http.MethodGet, "/v1/movies", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP status code 200 OK")
	assert.Contains(t, rr.Body.String(), "Inception", "Response body should contain 'Inception'")
}
