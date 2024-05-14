package main

import (
	"awesomeProject2/internal/data"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testServer struct {
	*httptest.Server
}

func newTestApp() *application {
	var db *pgxpool.Pool = nil

	mockMovies := data.MockMovieStorage{
		DB: db,
		InsertFunc: func(movie *data.Movie) error {
			movie.ID = 1
			movie.CreatedAt = time.Now()
			movie.Version = 1
			return nil
		},
		GetFunc: func(id int64) (*data.Movie, error) {
			if id == 1 {
				return &data.Movie{
					ID:        1,
					Title:     "Inception",
					Year:      2000,
					Runtime:   48,
					Genres:    []string{"action", "sci-fi"},
					CreatedAt: time.Now(),
					Version:   1,
				}, nil
			}
			return nil, errors.New("no record found")
		},
		GetAllFunc: func(title string, genres []string, filters data.Filters) ([]*data.Movie, data.Metadata, error) {
			filteredMovies := []*data.Movie{
				{ID: 1, Title: "Inception", Year: 2000, Runtime: 48, Genres: []string{"action", "sci-fi"}},
			}
			metadata := data.Metadata{TotalRecords: 1, FirstPage: 1, PageSize: 10}
			return filteredMovies, metadata, nil
		},
		UpdateFunc: func(movie *data.Movie) error {
			movie.Version++
			return nil
		},
		DeleteFunc: func(id int64) error {
			if id == 1 {
				return nil
			}
			return errors.New("movie not found")
		},
	}
	app := new(application)
	cfg := config{env: "testing"}
	app.config = cfg

	app.storage.Movies = &mockMovies
	return app
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := rs.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
