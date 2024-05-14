package main

import (
	"awesomeProject2/internal/data"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateMovieU(t *testing.T) {
	mockStorage := &data.MockMovieStorage{
		InsertFunc: func(movie *data.Movie) error {
			if movie.Title == "" {
				return errors.New("title is required")
			}
			if movie.Year == 0 {
				return errors.New("year is required")
			}
			movie.ID = 42
			return nil
		},
	}

	app := new(application)
	app.storage.Movies = mockStorage
	newMovie := &data.Movie{Title: "Inception", Year: 2000}
	err := app.storage.Movies.Insert(newMovie)

	assert.Nil(t, err)
	assert.NotEqual(t, 0, newMovie.ID, "New movie should have a non-zero ID after insertion")
}

func TestGetMovieByID(t *testing.T) {
	mockStorage := &data.MockMovieStorage{
		GetFunc: func(id int64) (*data.Movie, error) {
			if id != 1 {
				return nil, errors.New("movie not found")
			}
			return &data.Movie{ID: 1, Title: "Inception"}, nil
		},
	}

	app := new(application)
	app.storage.Movies = mockStorage

	movie, err := app.storage.Movies.Get(1)

	assert.Nil(t, err)
	assert.NotNil(t, movie, "Movie should not be nil")
	assert.Equal(t, "Inception", movie.Title, "The movie title should match")
}

func TestUpdateMovieU(t *testing.T) {
	mockStorage := &data.MockMovieStorage{
		UpdateFunc: func(movie *data.Movie) error {
			if movie.ID != 1 {
				return errors.New("movie not found")
			}
			movie.Version++
			return nil
		},
	}

	app := new(application)
	app.storage = data.Storage{Movies: mockStorage}

	movieToUpdate := &data.Movie{ID: 1, Title: "Interstellar", Version: 1}
	err := app.storage.Movies.Update(movieToUpdate)

	assert.Nil(t, err, "Expected no error on movie update")
	assert.Equal(t, int32(2), movieToUpdate.Version, "Version should be incremented")
}

func TestDeleteMovieUnit(t *testing.T) {
	mockStorage := &data.MockMovieStorage{
		DeleteFunc: func(id int64) error {
			if id != 1 {
				return errors.New("movie not found")
			}
			return nil
		},
	}

	app := new(application)
	app.storage.Movies = mockStorage

	err := app.storage.Movies.Delete(1)
	assert.Nil(t, err)

	err = app.storage.Movies.Delete(42)
	assert.NotNil(t, err, "Deleting a non-existent movie should result in an error")
}

func TestListMoviesU(t *testing.T) {
	mockStorage := &data.MockMovieStorage{
		GetAllFunc: func(title string, genres []string, filters data.Filters) ([]*data.Movie, data.Metadata, error) {
			movies := []*data.Movie{
				{ID: 1, Title: "Inception", Year: 2010, Runtime: 148, Genres: []string{"action", "sci-fi"}},
			}
			metadata := data.Metadata{TotalRecords: 1, FirstPage: 1, PageSize: 10}
			return movies, metadata, nil
		},
	}

	app := new(application)
	app.storage = data.Storage{Movies: mockStorage}

	title := ""
	genres := []string{}
	filters := data.Filters{}
	movies, metadata, err := app.storage.Movies.GetAll(title, genres, filters)

	assert.Nil(t, err, "Expected no error when fetching all movies")
	assert.Len(t, movies, 1, "Expected one movie to be returned")
	assert.Equal(t, "Inception", movies[0].Title, "Expected the title to be 'Inception'")
	assert.Equal(t, 1, metadata.TotalRecords, "Expected one record in metadata")
}

func TestHealthCheckU(t *testing.T) {
	app := newTestApp()
	router := app.routes()
	req, _ := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "testing")
}
